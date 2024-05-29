package api

import (
	"context"
	"fmt"
	"log"

	"github.com/H-BF/gw/internal/client/SecGroup"
	"github.com/H-BF/gw/pkg/authprovider"
	"github.com/H-BF/protos/pkg/api/sgroups"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"

	"connectrpc.com/connect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const userIDHeaderKey = "userId"

type SecGroupService struct {
	sgroupsconnect.UnimplementedSecGroupServiceHandler

	authPlugin authprovider.AuthProvider
	gwClient   sgroupsconnect.SecGroupServiceClient
}

func NewSecGroupService(authPlugin authprovider.AuthProvider) sgroupsconnect.SecGroupServiceHandler {
	gwClient := SecGroup.NewClient("http://localhost:9000")

	return &SecGroupService{
		authPlugin: authPlugin,
		gwClient:   gwClient,
	}
}

func (s SecGroupService) checkPermissions(ctx context.Context, reqTuples RTuples) error {
	for _, authReq := range reqTuples {
		isAuth, err := s.authPlugin.Authorize(ctx, authReq[0], authReq[1], authReq[2])
		if err != nil {
			log.Println(err)

			return status.Errorf(codes.Internal, "error during auth attempt: %v", err)
		}

		if !isAuth {
			log.Println(err)

			return status.Errorf(
				codes.PermissionDenied,
				"user %s does not have access or action permision to the %s resource, action - %s",
				authReq[0], authReq[1], authReq[2],
			)
		}
	}

	return nil
}

func (s SecGroupService) Sync(
	ctx context.Context,
	c *connect.Request[sgroups.SyncReq],
) (*connect.Response[emptypb.Empty], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromSync(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var objsToCreate []string

	for _, authReq := range tt {
		authResp, err := s.authPlugin.AuthorizeIfExist(ctx, authReq[0], authReq[1], authReq[2])
		if err != nil {
			log.Println(err)

			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}

		if !authResp.Authorized {
			err = fmt.Errorf(
				"user %s does not have access or action permision to the %s resource, action - %s",
				authReq[0], authReq[1], authReq[2],
			)

			log.Println(err)

			return nil, status.Error(
				codes.PermissionDenied,
				err.Error(),
			)
		}

		if !authResp.Exist {
			objsToCreate = append(objsToCreate, authReq[1])
		}
	}

	sgroupsResp, err := s.gwClient.Sync(ctx, c)
	if err == nil {
		// TODO: на подумать - как можно сделать изменения в группировках используя одну функцию?
		switch c.Msg.SyncOp {
		case sgroups.SyncReq_Delete:
			// todo: handle error
			if err = s.authPlugin.RemoveResourcesFromGroup(ctx, sub, tt.GetObjs()...); err != nil {
				log.Println(err)
			}
		case sgroups.SyncReq_Upsert:
			// todo: handle error
			if err = s.authPlugin.AddResourcesToGroup(ctx, sub, objsToCreate...); err != nil {
				log.Println(err)
			}
		}
	} else {
		log.Println(err)
	}

	return sgroupsResp, err
}

func (s SecGroupService) ListNetworks(
	ctx context.Context,
	c *connect.Request[sgroups.ListNetworksReq],
) (*connect.Response[sgroups.ListNetworksResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromListNetworks(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.ListNetworks(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) ListSecurityGroups(
	ctx context.Context,
	c *connect.Request[sgroups.ListSecurityGroupsReq],
) (*connect.Response[sgroups.ListSecurityGroupsResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromListSecurityGroups(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.ListSecurityGroups(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) GetRules(
	ctx context.Context,
	c *connect.Request[sgroups.GetRulesReq],
) (*connect.Response[sgroups.RulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromGetRules(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.GetRules(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) FindRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindRulesReq],
) (*connect.Response[sgroups.RulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindRules(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.FindRules(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) FindFqdnRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindFqdnRulesReq],
) (*connect.Response[sgroups.FqdnRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindFqdnRules(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.FindFqdnRules(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) FindCidrSgRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindCidrSgRulesReq],
) (*connect.Response[sgroups.CidrSgRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindCidrSgRules(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.FindCidrSgRules(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (s SecGroupService) FindSgSgRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindSgSgRulesReq],
) (*connect.Response[sgroups.SgSgRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindSgSgRules(c.Msg, sub); err != nil {
		log.Println(err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		log.Println(err)

		return nil, err
	}

	res, err := s.gwClient.FindSgSgRules(ctx, c)
	if err != nil {
		log.Println(err)
	}

	return res, err
}

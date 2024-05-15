package api

import (
	"context"
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
			return err
		}

		if !isAuth {
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromSync(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		objsToCreate []string
		objsToDelete []string
	)

	for _, authReq := range tt {
		authResp, err := s.authPlugin.AuthorizeIfExist(ctx, authReq[0], authReq[1], authReq[2])
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}

		if !authResp.Authorized {
			return nil, status.Errorf(
				codes.PermissionDenied,
				"user %s does not have access or action permision to the %s resource, action - %s",
				authReq[0], authReq[1], authReq[2],
			)
		}

		if c.Msg.SyncOp == sgroups.SyncReq_Delete {
			objsToDelete = append(objsToDelete, authReq[1])
			continue
		}

		if !authResp.Exist {
			objsToCreate = append(objsToCreate, authReq[1])
		}
	}

	sgroupsResp, err := s.gwClient.Sync(ctx, c)
	if err == nil {
		// todo: handle error
		if c.Msg.SyncOp == sgroups.SyncReq_Delete {
			if err = s.authPlugin.RemoveResourcesFromGroup(ctx, sub, objsToDelete...); err != nil {
				log.Println(err)
			}
		}

		// todo: handle error
		if len(objsToCreate) > 0 {
			if err = s.authPlugin.AddResourcesToGroup(ctx, sub, objsToCreate...); err != nil {
				log.Println(err)
			}
		}
	}

	return sgroupsResp, err
}

func (s SecGroupService) ListNetworks(
	ctx context.Context,
	c *connect.Request[sgroups.ListNetworksReq],
) (*connect.Response[sgroups.ListNetworksResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromListNetworks(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.ListNetworks(ctx, c)
}

func (s SecGroupService) ListSecurityGroups(
	ctx context.Context,
	c *connect.Request[sgroups.ListSecurityGroupsReq],
) (*connect.Response[sgroups.ListSecurityGroupsResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromListSecurityGroups(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.ListSecurityGroups(ctx, c)
}

func (s SecGroupService) GetRules(
	ctx context.Context,
	c *connect.Request[sgroups.GetRulesReq],
) (*connect.Response[sgroups.RulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromGetRules(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.GetRules(ctx, c)
}

func (s SecGroupService) FindRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindRulesReq],
) (*connect.Response[sgroups.RulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindRules(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.FindRules(ctx, c)
}

func (s SecGroupService) FindFqdnRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindFqdnRulesReq],
) (*connect.Response[sgroups.FqdnRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindFqdnRules(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.FindFqdnRules(ctx, c)
}

func (s SecGroupService) FindCidrSgRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindCidrSgRulesReq],
) (*connect.Response[sgroups.CidrSgRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindCidrSgRules(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.FindCidrSgRules(ctx, c)
}

func (s SecGroupService) FindSgSgRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindSgSgRulesReq],
) (*connect.Response[sgroups.SgSgRulesResp], error) {
	sub, err := extractSub(c)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var tt RTuples
	if err := tt.FromFindSgSgRules(c.Msg, sub); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.checkPermissions(ctx, tt); err != nil {
		return nil, err
	}

	return s.gwClient.FindSgSgRules(ctx, c)
}

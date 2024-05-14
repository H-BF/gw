package api

import (
	"context"

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

	var toCreate RTuples

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

		if !authResp.Exist {
			toCreate = append(toCreate, authReq)
		}
	}

	sgroupsResp, err := s.gwClient.Sync(ctx, c)
	if err == nil {
		// TODO: add objs to groups in polices
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

func (s SecGroupService) SyncStatus(
	ctx context.Context,
	c *connect.Request[emptypb.Empty],
) (*connect.Response[sgroups.SyncStatusResp], error) {
	return nil, status.Error(codes.Unimplemented, "method SyncStatus not implemented")
}

func (s SecGroupService) SyncStatuses(
	ctx context.Context,
	c *connect.Request[emptypb.Empty],
	c2 *connect.ServerStream[sgroups.SyncStatusResp],
) error {
	return status.Error(codes.Unimplemented, "method SyncStatuses not implemented")
}

func (s SecGroupService) GetSgSubnets(
	ctx context.Context,
	c *connect.Request[sgroups.GetSgSubnetsReq],
) (*connect.Response[sgroups.GetSgSubnetsResp], error) {
	return nil, status.Error(codes.Unimplemented, "method GetSgSubnets not implemented")
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

func (s SecGroupService) FindSgIcmpRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindSgIcmpRulesReq],
) (*connect.Response[sgroups.SgIcmpRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindSgIcmpRules not implemented")
}

func (s SecGroupService) FindSgSgIcmpRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindSgSgIcmpRulesReq],
) (*connect.Response[sgroups.SgSgIcmpRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindSgSgIcmpRules not implemented")
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

func (s SecGroupService) FindIESgSgIcmpRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindIESgSgIcmpRulesReq],
) (*connect.Response[sgroups.IESgSgIcmpRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindIESgSgIcmpRules not implemented")
}

func (s SecGroupService) FindCidrSgIcmpRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindCidrSgIcmpRulesReq],
) (*connect.Response[sgroups.CidrSgIcmpRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindCidrSgIcmpRules not implemented")
}

func (s SecGroupService) GetSecGroupForAddress(
	ctx context.Context,
	c *connect.Request[sgroups.GetSecGroupForAddressReq],
) (*connect.Response[sgroups.SecGroup], error) {
	return nil, status.Error(codes.Unimplemented, "method GetSecGroupForAddress not implemented")
}

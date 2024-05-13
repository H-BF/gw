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

func (s SecGroupService) checkPermission(ctx context.Context, sub, obj, act string) error {
	isAuth, err := s.authPlugin.CheckPermission(ctx, sub, obj, act)
	if err != nil {
		return err
	}

	if !isAuth {
		return status.Errorf(
			codes.PermissionDenied,
			"user %s does not have access or action permision to the %s resource, action - %s",
			sub, obj, act,
		)
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
	for _, t := range tt {
		if err := s.checkPermission(ctx, t[0], t[1], t[2]); err != nil {
			return nil, err
		}
	}

	return s.gwClient.Sync(ctx, c)
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
	for _, t := range tt {
		if err := s.checkPermission(ctx, t[0], t[1], t[2]); err != nil {
			return nil, err
		}
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
	for _, t := range tt {
		if err := s.checkPermission(ctx, t[0], t[1], t[2]); err != nil {
			return nil, err
		}
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
	for _, t := range tt {
		if err := s.checkPermission(ctx, t[0], t[1], t[2]); err != nil {
			return nil, err
		}
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
	return nil, status.Error(codes.Unimplemented, "method FindRules not implemented")
}

func (s SecGroupService) FindFqdnRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindFqdnRulesReq],
) (*connect.Response[sgroups.FqdnRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindFqdnRules not implemented")
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
	return nil, status.Error(codes.Unimplemented, "method FindCidrSgRules not implemented")
}

func (s SecGroupService) FindSgSgRules(
	ctx context.Context,
	c *connect.Request[sgroups.FindSgSgRulesReq],
) (*connect.Response[sgroups.SgSgRulesResp], error) {
	return nil, status.Error(codes.Unimplemented, "method FindSgSgRules not implemented")
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

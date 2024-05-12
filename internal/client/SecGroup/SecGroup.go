package SecGroup

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type secGroupClient struct {
	client sgroupsconnect.SecGroupServiceClient
}

func NewClient(addr string) sgroupsconnect.SecGroupServiceClient {
	// client for transferring data via http2 to sgroups grpc server
	httpClient := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	client := sgroupsconnect.NewSecGroupServiceClient(
		httpClient,
		addr,
		connect.WithGRPC(),
	)

	return secGroupClient{client: client}
}

func (s secGroupClient) Sync(
	ctx context.Context,
	c *connect.Request[sgroups.SyncReq],
) (*connect.Response[emptypb.Empty], error) {
	return s.client.Sync(ctx, c)
}

func (s secGroupClient) SyncStatus(ctx context.Context, c *connect.Request[emptypb.Empty]) (*connect.Response[sgroups.SyncStatusResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) SyncStatuses(ctx context.Context, c *connect.Request[emptypb.Empty]) (*connect.ServerStreamForClient[sgroups.SyncStatusResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) ListNetworks(
	ctx context.Context,
	c *connect.Request[sgroups.ListNetworksReq],
) (*connect.Response[sgroups.ListNetworksResp], error) {
	return s.client.ListNetworks(ctx, c)
}

func (s secGroupClient) ListSecurityGroups(
	ctx context.Context,
	c *connect.Request[sgroups.ListSecurityGroupsReq],
) (*connect.Response[sgroups.ListSecurityGroupsResp], error) {
	return s.client.ListSecurityGroups(ctx, c)
}

func (s secGroupClient) GetSgSubnets(ctx context.Context, c *connect.Request[sgroups.GetSgSubnetsReq]) (*connect.Response[sgroups.GetSgSubnetsResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) GetRules(
	ctx context.Context,
	c *connect.Request[sgroups.GetRulesReq],
) (*connect.Response[sgroups.RulesResp], error) {
	return s.client.GetRules(ctx, c)
}

func (s secGroupClient) FindRules(ctx context.Context, c *connect.Request[sgroups.FindRulesReq]) (*connect.Response[sgroups.RulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindFqdnRules(ctx context.Context, c *connect.Request[sgroups.FindFqdnRulesReq]) (*connect.Response[sgroups.FqdnRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindSgIcmpRules(ctx context.Context, c *connect.Request[sgroups.FindSgIcmpRulesReq]) (*connect.Response[sgroups.SgIcmpRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindSgSgIcmpRules(ctx context.Context, c *connect.Request[sgroups.FindSgSgIcmpRulesReq]) (*connect.Response[sgroups.SgSgIcmpRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindCidrSgRules(ctx context.Context, c *connect.Request[sgroups.FindCidrSgRulesReq]) (*connect.Response[sgroups.CidrSgRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindSgSgRules(ctx context.Context, c *connect.Request[sgroups.FindSgSgRulesReq]) (*connect.Response[sgroups.SgSgRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindIESgSgIcmpRules(ctx context.Context, c *connect.Request[sgroups.FindIESgSgIcmpRulesReq]) (*connect.Response[sgroups.IESgSgIcmpRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) FindCidrSgIcmpRules(ctx context.Context, c *connect.Request[sgroups.FindCidrSgIcmpRulesReq]) (*connect.Response[sgroups.CidrSgIcmpRulesResp], error) {
	//TODO implement me
	panic("implement me")
}

func (s secGroupClient) GetSecGroupForAddress(ctx context.Context, c *connect.Request[sgroups.GetSecGroupForAddressReq]) (*connect.Response[sgroups.SecGroup], error) {
	//TODO implement me
	panic("implement me")
}

package api

import (
	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups"
)

type SecGroupReq interface {
	sgroups.ListNetworksReq | sgroups.ListSecurityGroupsReq | sgroups.SyncReq | sgroups.GetRulesReq |
		sgroups.FindRulesReq | sgroups.FindFqdnRulesReq | sgroups.FindCidrSgRulesReq | sgroups.FindSgSgRulesReq
}

func extractSub[T SecGroupReq](c *connect.Request[T]) (string, error) {
	return c.Header().Get(userIDHeaderKey), nil
}

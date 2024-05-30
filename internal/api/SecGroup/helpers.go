package api

import (
	"connectrpc.com/connect"
	"errors"
	"github.com/H-BF/protos/pkg/api/sgroups"
	"strings"
)

type SecGroupReq interface {
	sgroups.ListNetworksReq | sgroups.ListSecurityGroupsReq | sgroups.SyncReq | sgroups.GetRulesReq |
		sgroups.FindRulesReq | sgroups.FindFqdnRulesReq | sgroups.FindCidrSgRulesReq | sgroups.FindSgSgRulesReq
}

var (
	errEmptyAuthnInfo = errors.New("empty authentication info")
)

func extractSub[T SecGroupReq](c *connect.Request[T]) (sub string, err error) {
	sub = strings.TrimSpace(c.Header().Get(userIDHeaderKey))
	if sub == "" {
		err = errEmptyAuthnInfo
	}
	return sub, err
}

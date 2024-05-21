package api

import (
	"errors"
	"strings"

	"github.com/H-BF/gw/internal/authprovider/consts"

	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups"
)

type SecGroupReq interface {
	sgroups.ListNetworksReq | sgroups.ListSecurityGroupsReq | sgroups.SyncReq | sgroups.GetRulesReq |
		sgroups.FindRulesReq | sgroups.FindFqdnRulesReq | sgroups.FindCidrSgRulesReq | sgroups.FindSgSgRulesReq
}

func extractSub[T SecGroupReq](c *connect.Request[T]) (string, error) {
	userId := c.Header().Get(userIDHeaderKey)
	if strings.TrimSpace(userId) == "" {
		return "", errors.New("user id connot be empty")
	}

	return userId, nil
}

func extractAct(req *sgroups.SyncReq) (string, error) {
	switch req.GetSyncOp() {
	case sgroups.SyncReq_Upsert, sgroups.SyncReq_Delete:
		return consts.EditAction, nil
	default:
		return "", errUnsupportedSyncOp
	}
}

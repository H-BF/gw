package api

import (
	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups"
	"github.com/H-BF/sgroups-k8s-adapter/internal/authprovider"
)

func getSyncResourceByRequest(request *connect.Request[sgroups.SyncReq]) string {
	if request.Msg.GetNetworks() != nil {
		return authprovider.NETWORK
	}

	if request.Msg.GetFqdnRules() != nil {
		return authprovider.S2F
	}

	if request.Msg.GetGroups() != nil {
		return authprovider.SECURITY_GROUP
	}

	return ""
}

func getActionBySyncOp(syncOp string) string {
	switch syncOp {
	case "Upsert", "Delete":
		return authprovider.EditAction
	default:
		return ""
	}
}

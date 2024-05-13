package api

import (
	"errors"

	ap "github.com/H-BF/gw/internal/authprovider"

	"github.com/H-BF/protos/pkg/api/sgroups"
)

type RTuples [][3]string

var (
	errUnsupportedSyncOp = errors.New("unsupported SyncOp: only Upsert, Delete allowed")
)

func extractAct(req *sgroups.SyncReq) (string, error) {
	switch req.GetSyncOp() {
	case sgroups.SyncReq_Upsert, sgroups.SyncReq_Delete:
		return ap.EditAction, nil
	default:
		return "", errUnsupportedSyncOp
	}
}

func (t *RTuples) FromSync(req *sgroups.SyncReq, sub string) error {
	var (
		act string
		err error
	)
	if act, err = extractAct(req); err != nil {
		return err
	}
	switch req.GetSubject().(type) {
	case *sgroups.SyncReq_Networks:
		for _, nw := range req.GetNetworks().GetNetworks() {
			*t = append(*t, [3]string{sub, nw.GetName(), act})
		}
	case *sgroups.SyncReq_Groups:
		for _, sg := range req.GetGroups().GetGroups() {
			*t = append(*t, [3]string{sub, sg.GetName(), act})
		}
	case *sgroups.SyncReq_FqdnRules:
		for _, s2f := range req.GetFqdnRules().GetRules() {
			*t = append(*t, [3]string{sub, s2f.GetSgFrom(), act})
		}
	default:
		return errors.New("unsupported sync subject")
	}
	return nil
}

func (t *RTuples) FromListNetworks(req *sgroups.ListNetworksReq, sub string) error {
	for _, obj := range req.GetNeteworkNames() {
		*t = append(*t, [3]string{sub, obj, ap.ReadAction})
	}
	return nil
}

func (t *RTuples) FromListSecurityGroups(req *sgroups.ListSecurityGroupsReq, sub string) error {
	for _, obj := range req.GetSgNames() {
		*t = append(*t, [3]string{sub, obj, ap.ReadAction})
	}
	return nil
}

func (t *RTuples) FromGetRules(req *sgroups.GetRulesReq, sub string) error {
	*t = append(*t, [3]string{sub, req.GetSgFrom(), ap.ReadAction})
	*t = append(*t, [3]string{sub, req.GetSgTo(), ap.ReadAction})
	return nil
}

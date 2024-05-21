package api

import (
	"errors"

	"github.com/H-BF/gw/internal/api/SecGroup/resnaming"
	"github.com/H-BF/gw/internal/authprovider/consts"

	"github.com/H-BF/protos/pkg/api/sgroups"
)

type RTuples [][3]string

var (
	errUnsupportedSyncOp = errors.New("unsupported SyncOp: only Upsert, Delete allowed")
)

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
		for _, rule := range req.GetFqdnRules().GetRules() {
			*t = append(*t, [3]string{sub, resnaming.RuleName(rule), act})
			*t = append(*t, [3]string{sub, rule.GetSgFrom(), consts.ReferenceAction})
		}
	case *sgroups.SyncReq_SgRules:
		for _, rule := range req.GetSgRules().GetRules() {
			*t = append(*t, [3]string{sub, resnaming.RuleName(rule), act})
			*t = append(*t, [3]string{sub, rule.GetSgFrom(), consts.ReferenceAction})
			*t = append(*t, [3]string{sub, rule.GetSgTo(), consts.ReferenceAction})
		}
	case *sgroups.SyncReq_CidrSgRules:
		for _, rule := range req.GetCidrSgRules().GetRules() {
			*t = append(*t, [3]string{sub, resnaming.RuleName(rule), act})
			*t = append(*t, [3]string{sub, rule.GetSG(), consts.ReferenceAction})
		}
	case *sgroups.SyncReq_SgSgRules:
		for _, rule := range req.GetSgSgRules().GetRules() {
			*t = append(*t, [3]string{sub, resnaming.RuleName(rule), act})
			*t = append(*t, [3]string{sub, rule.GetSgLocal(), consts.ReferenceAction})
		}
	default:
		return errors.New("unsupported sync subject")
	}
	return nil
}

func (t *RTuples) FromListNetworks(req *sgroups.ListNetworksReq, sub string) error {
	if len(req.GetNeteworkNames()) == 0 {
		// TODO: если пришел запрос с пустым массивом имен нетворков то сгрупс отдаст все существующие нетворки
		// поэтому нужно проверить у пользователя доступ на чтение нетворков
		// то же самое для всех List* и Find* запросов - их тоже нужно переделать
		*t = append(*t, [3]string{sub, "all-networks", consts.ReadAction})
	}
	for _, obj := range req.GetNeteworkNames() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) FromListSecurityGroups(req *sgroups.ListSecurityGroupsReq, sub string) error {
	for _, obj := range req.GetSgNames() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) FromGetRules(req *sgroups.GetRulesReq, sub string) error {
	*t = append(*t, [3]string{sub, req.GetSgFrom(), consts.ReadAction})
	*t = append(*t, [3]string{sub, req.GetSgTo(), consts.ReadAction})
	return nil
}

func (t *RTuples) FromFindRules(req *sgroups.FindRulesReq, sub string) error {
	for _, obj := range req.GetSgFrom() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	for _, obj := range req.GetSgTo() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) FromFindFqdnRules(req *sgroups.FindFqdnRulesReq, sub string) error {
	for _, obj := range req.GetSgFrom() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) FromFindCidrSgRules(req *sgroups.FindCidrSgRulesReq, sub string) error {
	for _, obj := range req.GetSg() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) FromFindSgSgRules(req *sgroups.FindSgSgRulesReq, sub string) error {
	for _, obj := range req.GetSgLocal() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	for _, obj := range req.GetSg() {
		*t = append(*t, [3]string{sub, obj, consts.ReadAction})
	}
	return nil
}

func (t *RTuples) GetObjs() (res []string) {
	for _, tuple := range *t {
		res = append(res, tuple[1])
	}
	return res
}

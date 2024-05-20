package resnaming

import (
	"fmt"
	"strings"

	"github.com/H-BF/protos/pkg/api/sgroups"
)

func RuleName(rule interface{}) string {
	switch rule := rule.(type) {
	case *sgroups.FqdnRule:
		return fmt.Sprintf(
			"%s:sg(%s)fqdn(%s)",
			strings.ToLower(rule.GetTransport().String()),
			strings.ToLower(rule.GetSgFrom()),
			strings.ToLower(rule.GetFQDN()),
		)
	case *sgroups.Rule:
		return fmt.Sprintf(
			"%s:sg(%s)sg(%s)",
			strings.ToLower(rule.GetTransport().String()),
			strings.ToLower(rule.GetSgFrom()),
			strings.ToLower(rule.GetSgTo()),
		)
	case *sgroups.CidrSgRule:
		return fmt.Sprintf(
			"%s:cidr(%s)sg(%s)%s",
			strings.ToLower(rule.GetTransport().String()),
			strings.ToLower(rule.GetCIDR()),
			strings.ToLower(rule.GetSG()),
			strings.ToLower(rule.GetTraffic().String()),
		)
	case *sgroups.SgSgRule:
		return fmt.Sprintf(
			"%s:sg-local(%s)sg(%s)%s",
			strings.ToLower(rule.GetTransport().String()),
			strings.ToLower(rule.GetSgLocal()),
			strings.ToLower(rule.GetSg()),
			strings.ToLower(rule.GetTraffic().String()),
		)
	default:
		return ""
	}
}

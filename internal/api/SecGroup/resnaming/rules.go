package resnaming

import (
	"fmt"
	"strings"

	"github.com/H-BF/protos/pkg/api/sgroups"
)

func RuleName(rule interface{}) string {
	switch rule.(type) {
	case *sgroups.FqdnRule:
		fqdnRule := rule.(*sgroups.FqdnRule)
		return fmt.Sprintf(
			"%s:sg(%s)fqdn(%s)",
			strings.ToLower(fqdnRule.GetTransport().String()),
			strings.ToLower(fqdnRule.GetSgFrom()),
			strings.ToLower(fqdnRule.GetFQDN()),
		)
	case *sgroups.Rule:
		sgRule := rule.(*sgroups.Rule)
		return fmt.Sprintf(
			"%s:sg(%s)sg(%s)",
			strings.ToLower(sgRule.GetTransport().String()),
			strings.ToLower(sgRule.GetSgFrom()),
			strings.ToLower(sgRule.GetSgTo()),
		)
	case *sgroups.CidrSgRule:
		cidrRule := rule.(*sgroups.CidrSgRule)
		return fmt.Sprintf(
			"%s:cidr(%s)sg(%s)%s",
			strings.ToLower(cidrRule.GetTransport().String()),
			strings.ToLower(cidrRule.GetCIDR()),
			strings.ToLower(cidrRule.GetSG()),
			strings.ToLower(cidrRule.GetTraffic().String()),
		)
	case *sgroups.SgSgRule:
		sgSgRule := rule.(*sgroups.SgSgRule)
		return fmt.Sprintf(
			"%s:sg-local(%s)sg(%s)%s",
			sgSgRule.GetTransport().String(),
			sgSgRule.GetSgLocal(),
			sgSgRule.GetSg(),
			sgSgRule.GetTraffic().String(),
		)
	default:
		return ""
	}
}

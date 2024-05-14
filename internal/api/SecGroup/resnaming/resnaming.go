package resnaming

import (
	"fmt"
	"strings"
)

type (
	FQDNRle struct {
		Transport string
		SgFrom    string
		FqdnTo    string
	}
	SqRules struct {
		Transport string
		SgFrom    string
		SgTo      string
	}
	CIDRRule struct {
		Transport string
		SgName    string
		CIDR      string
		Traffic   string
	}
)

func FQNDRuleString(fr FQDNRle) string {
	return fmt.Sprintf(
		"%s:sg(%s)fqdn(%s)",
		strings.ToLower(fr.Transport),
		strings.ToLower(fr.SgFrom),
		strings.ToLower(fr.FqdnTo),
	)
}

func SgRuleString(sr SqRules) string {
	return fmt.Sprintf(
		"%s:sg(%s)sg(%s)",
		strings.ToLower(sr.Transport),
		strings.ToLower(sr.SgFrom),
		strings.ToLower(sr.SgTo),
	)
}

func CIDRRUleString(cr CIDRRule) string {
	return fmt.Sprintf(
		"%s:cidr(%s)sg(%s)%s",
		strings.ToLower(cr.Transport),
		strings.ToLower(cr.CIDR),
		strings.ToLower(cr.SgName),
		strings.ToLower(cr.Traffic),
	)
}

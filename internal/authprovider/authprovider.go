package authprovider

import (
	"context"
	"errors"
	"strings"

	"github.com/H-BF/sgroups-k8s-adapter/pkg/authprovider"

	"github.com/casbin/casbin/v2"
)

type CasbinAuthProvider struct {
	enforcer *casbin.Enforcer
}

const ( // available actions of a role model
	ReadAction      = "read"
	EditAction      = "edit"
	ReferenceAction = "reference"
)

const ( // available role model resources
	NETWORK        = "Network"
	SECURITY_GROUP = "SecurityGroup"
	FQDN_S2F       = "FQDNS2F"
)

const (
	networkPrefix       = "nw-"
	securityGroupPrefix = "sg-"
	fqdnRulePrefix      = "fqdn-"
)

func NewCasbinAuthProvider(modelPath, policyPath string) (authprovider.AuthProvider, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return &CasbinAuthProvider{enforcer}, nil
}

func (c CasbinAuthProvider) CheckPermission(_ context.Context, sub, obj, act string) (bool, error) {
	if _, err := c.addResourceToNamedGroup(obj); err != nil {
		return false, err
	}

	return c.enforcer.Enforce(sub, obj, act)
}

func (c CasbinAuthProvider) addResourceToNamedGroup(resourceName string) (bool, error) {
	var group string

	switch {
	case strings.HasPrefix(resourceName, networkPrefix):
		group = NETWORK
	case strings.HasPrefix(resourceName, securityGroupPrefix):
		group = SECURITY_GROUP
	case strings.HasPrefix(resourceName, fqdnRulePrefix):
		group = FQDN_S2F
	default:
		return false, errors.New("unknown resource type")
	}

	// if the resource has already been created in the group,
	// then a new entry will not be written.
	isAdded, err := c.enforcer.AddNamedGroupingPolicy("g2", group, resourceName)
	if err != nil {
		return false, err
	}

	if err = c.enforcer.SavePolicy(); err != nil {
		return false, err
	}

	return isAdded, nil
}

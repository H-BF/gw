package authprovider

import (
	"context"
	"github.com/H-BF/sgroups-k8s-adapter/pkg/authprovider"
	"github.com/casbin/casbin/v2"
	"strings"
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
	S2F            = "S2F"
)

const (
	networkPrefix       = "nw-"
	securityGroupPrefix = "sg-"
	rulePrefix          = "s2f-"
)

func NewCasbinAuthProvider(modelPath, policyPath string) (authprovider.AuthProvider, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	adapter, err := newCasbinFileAdapter(policyPath, enforcer.GetAdapter())
	if err != nil {
		return nil, err
	}

	enforcer.SetAdapter(adapter)
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &CasbinAuthProvider{enforcer}, nil
}

func (c CasbinAuthProvider) CheckPermission(_ context.Context, sub, obj, act string) (bool, error) {
	allRoles := c.enforcer.GetAllRoles()
	for _, role := range allRoles {
		switch {
		case role == obj:
			return c.enforcer.Enforce(sub, obj, act)
		case strings.HasPrefix(role, networkPrefix):
			return c.enforcer.AddPolicy([]string{"g2", NETWORK, obj})
		case strings.HasPrefix(role, securityGroupPrefix):
			return c.enforcer.AddPolicy([]string{"g2", SECURITY_GROUP, obj})
		case strings.HasPrefix(role, rulePrefix):
			return c.enforcer.AddPolicy([]string{"g2", S2F, obj})
		}
	}

	return c.enforcer.Enforce(sub, obj, act)
}

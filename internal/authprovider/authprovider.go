package authprovider

import (
	"context"
	"github.com/H-BF/gw/pkg/authprovider"

	"github.com/casbin/casbin/v2"
)

type CasbinAuthProvider struct {
	enforcer *casbin.Enforcer
}

var _ authprovider.AuthProvider = (*CasbinAuthProvider)(nil)

const ( // available actions of a role model
	ReadAction      = "read"
	EditAction      = "edit"
	ReferenceAction = "reference"
)

const (
	adminRole = "admin"
	ownerRole = "owner"
)

const (
	subGroupSuffix = "-res"

	G2 = "g2"
)

func NewCasbinAuthProvider(modelPath, policyPath string) (authprovider.AuthProvider, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return &CasbinAuthProvider{enforcer}, nil
}

// Authorize implements authprovider.AuthProvider
func (c CasbinAuthProvider) Authorize(_ context.Context, sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

// AuthorizeIfExist implements authprovider.AuthProvider
func (c CasbinAuthProvider) AuthorizeIfExist(_ context.Context, sub, obj, act string) (authprovider.AuthWithExistResp, error) {
	if exists := c.enforcer.HasNamedGroupingPolicy("g2", sub+subGroupSuffix, obj); !exists {
		// if `obj` not added to group resource should be authorized and added to group after succeeded request
		return authprovider.AuthWithExistResp{
			Exist:      false,
			Authorized: true,
		}, nil
	}

	authorized, err := c.enforcer.Enforce(sub, obj, act)
	return authprovider.AuthWithExistResp{
		Exist:      false,
		Authorized: authorized,
	}, err
}

// AddResourcesToGroup implements authprovider.AuthProvider
func (c CasbinAuthProvider) AddResourcesToGroup(_ context.Context, sub string, objs ...string) error {
	for _, obj := range objs {
		if _, err := c.enforcer.AddNamedGroupingPolicy(G2, sub+subGroupSuffix, obj); err != nil {
			return err
		}
	}
	return c.enforcer.SavePolicy()
}

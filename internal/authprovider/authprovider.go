package authprovider

import (
	"context"
	"github.com/H-BF/gw/pkg/authprovider"

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

func NewCasbinAuthProvider(modelPath, policyPath string) (authprovider.AuthProvider, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return &CasbinAuthProvider{enforcer}, nil
}

func (c CasbinAuthProvider) CheckPermission(_ context.Context, sub, obj, act string) (bool, error) {
	if _, err := c.addResourceToNamedGroup(sub, obj); err != nil {
		return false, err
	}

	return c.enforcer.Enforce(sub, obj, act)
}

func (c CasbinAuthProvider) addResourceToNamedGroup(userId, resourceName string) (bool, error) {

	// if the resource has already been created in the group,
	// then a new entry will not be written.
	isAdded, err := c.enforcer.AddNamedGroupingPolicy("g2", userId+"-res", resourceName)
	if err != nil {
		return false, err
	}

	if err = c.enforcer.SavePolicy(); err != nil {
		return false, err
	}

	return isAdded, nil
}

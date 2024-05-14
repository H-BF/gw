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

const (
	adminRole = "admin"
	ownerRole = "owner"
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
	// for now, for each user who tries to create a new resource,
	// we will select the user role
	// todo: does it make sense t give admin the owner role???
	if err := c.addRoleForUser(userId, "owner"); err != nil {
		return false, err
	}

	// if the resource has already been created in the group,
	// then a new entry will not be written.
	isAdded, err := c.enforcer.AddNamedGroupingPolicy("g2", userId+"-res", resourceName)
	if err != nil {
		return false, err
	}

	// todo: save police only if all auth steps are successful
	if err = c.enforcer.SavePolicy(); err != nil {
		return false, err
	}

	return isAdded, nil
}

func (c CasbinAuthProvider) addRoleForUser(sub, role string) error {
	hasRole, err := c.enforcer.HasRoleForUser(sub, role)
	if err != nil {
		return err
	}

	if !hasRole {
		if _, err = c.enforcer.AddPolicy(sub, sub+"-res", ownerRole); err != nil {
			return err
		}
	}

	return nil
}

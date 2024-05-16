package authprovider

import (
	"context"
	"fmt"
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
	// todo: make a more beautiful solution so that there is no coping code
	if !c.subExists(sub) {
		return false, fmt.Errorf("you cannot add a resource to an existing user - %s", sub) // сообщение в ошибке сбивает с толку - не делай так
	}
	// TODO: ^^^^^^ нам не нужно самостоятельно лезть в касбин для проверки авторизации - это сделает enforcer.Enforce
	// УДАЛИТЬ!!!

	return c.enforcer.Enforce(sub, obj, act)
}

// AuthorizeIfExist implements authprovider.AuthProvider
func (c CasbinAuthProvider) AuthorizeIfExist(_ context.Context, sub, obj, act string) (authprovider.AuthWithExistResp, error) {
	authorized, reasons, err := c.enforcer.EnforceEx(sub, obj, act)
	if !authorized && len(reasons) == 0 { // `obj` is new for casbin
		if exists := c.objExists(obj); !exists {
			// if `obj` not added to group resource should be authorized and added to group after succeeded request
			return authprovider.AuthWithExistResp{
				Exist:      false,
				Authorized: true,
			}, nil
		}
	}

	return authprovider.AuthWithExistResp{
		Exist:      true,
		Authorized: authorized,
	}, err
}

// AddResourcesToGroup implements authprovider.AuthProvider
func (c CasbinAuthProvider) AddResourcesToGroup(_ context.Context, sub string, objs ...string) error {
	if !c.subExists(sub) {
		return fmt.Errorf("you cannot add a resource to an existing user - %s", sub)
	}

	for _, obj := range objs {
		if _, err := c.enforcer.AddNamedGroupingPolicy(G2, sub+subGroupSuffix, obj); err != nil {
			return err
		}
	}
	return c.enforcer.SavePolicy()
}

// RemoveResourcesFromGroup implements authprovider.RemoveResourceFromGroup
func (c CasbinAuthProvider) RemoveResourcesFromGroup(_ context.Context, sub string, objs ...string) error {
	for _, obj := range objs {
		if _, err := c.enforcer.RemoveNamedGroupingPolicy(G2, sub+subGroupSuffix, obj); err != nil {
			return err
		}
	}

	return c.enforcer.SavePolicy()
}

func (c CasbinAuthProvider) objExists(obj string) bool {
	const (
		groupNameIndex = 0
		objIndex       = 1
	)
	policy := c.enforcer.GetFilteredNamedGroupingPolicy(G2, objIndex, obj)
	return len(policy) != 0
}

// TODO: удалить
func (c CasbinAuthProvider) subExists(sub string) bool {
	policy := c.enforcer.GetFilteredPolicy(0, sub)
	return len(policy) >= 1
}

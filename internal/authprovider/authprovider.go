package authprovider

import (
	"context"

	"github.com/H-BF/gw/pkg/authprovider"

	"github.com/casbin/casbin/v2"

	"github.com/pkg/errors"
)

type CasbinAuthProvider struct {
	enforcer *casbin.Enforcer
}

var _ authprovider.AuthProvider = (*CasbinAuthProvider)(nil)

const pkgApi = "authprovider."

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
	const api = pkgApi + "NewCasbinAuthProvider"

	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, errors.Wrap(err, api)
	}

	return &CasbinAuthProvider{enforcer}, nil
}

// Authorize implements authprovider.AuthProvider
func (c CasbinAuthProvider) Authorize(_ context.Context, sub, obj, act string) (bool, error) {
	const api = pkgApi + "CasbinAuthProvider.Authorize"

	// todo: make a more beautiful solution so that there is no coping code
	if !c.subExists(sub) {
		return false, errors.Wrap(errors.New("user does not exist in the system"), api)
	}
	// TODO: ^^^^^^ нам не нужно самостоятельно лезть в касбин для проверки авторизации - это сделает enforcer.Enforce
	// УДАЛИТЬ!!!

	isAuth, err := c.enforcer.Enforce(sub, obj, act)

	return isAuth, errors.Wrap(err, api)
}

// AuthorizeIfExist implements authprovider.AuthProvider
func (c CasbinAuthProvider) AuthorizeIfExist(_ context.Context, sub, obj, act string) (authprovider.AuthWithExistResp, error) {
	const api = pkgApi + "CasbinAuthProvider.AuthorizeIfExist"

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
	}, errors.Wrap(err, api)
}

// AddResourcesToGroup implements authprovider.AuthProvider
func (c CasbinAuthProvider) AddResourcesToGroup(_ context.Context, sub string, objs ...string) error {
	const api = pkgApi + "CasbinAuthProvider.AddResourcesToGroup"

	if !c.subExists(sub) {
		return errors.Wrap(errors.Errorf("you cannot add a resource to an existing user - %s", sub), api)
	}

	for _, obj := range objs {
		if _, err := c.enforcer.AddNamedGroupingPolicy(G2, sub+subGroupSuffix, obj); err != nil {
			return errors.Wrap(errors.Errorf("an error occurred during created of the %s resource: %v", obj, err), api)
		}
	}

	return errors.Wrap(c.enforcer.SavePolicy(), api)
}

// RemoveResourcesFromGroup implements authprovider.RemoveResourceFromGroup
func (c CasbinAuthProvider) RemoveResourcesFromGroup(_ context.Context, sub string, objs ...string) error {
	const api = pkgApi + "CasbinAuthProvider.RemoveResourcesFromGroup"

	for _, obj := range objs {
		if _, err := c.enforcer.RemoveNamedGroupingPolicy(G2, sub+subGroupSuffix, obj); err != nil {
			return errors.Wrap(errors.Errorf("an error occurred during deletion of the %s resource: %v", obj, err), api)
		}
	}

	return errors.Wrap(c.enforcer.SavePolicy(), api)
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

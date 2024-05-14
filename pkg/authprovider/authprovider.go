package authprovider

import (
	"context"
)

type AuthWithExistResp struct {
	Exist      bool
	Authorized bool
}

type AuthProvider interface {
	// Authorize does authorization for passed request tuple {sub, obj, act}
	Authorize(ctx context.Context, sub, obj, act string) (bool, error)

	// AuthorizeIfExist does authorization and returns flag for existence of `sub` in policies
	AuthorizeIfExist(ctx context.Context, sub, obj, act string) (AuthWithExistResp, error)

	// AddResourcesToGroup adds `objs` to resource group of `sub`
	AddResourcesToGroup(ctx context.Context, sub string, objs ...string) error
}

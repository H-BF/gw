package authprovider

import "context"

type AuthProvider interface {
	CheckPermission(ctx context.Context, sub, obj, act string) (bool, error)
}

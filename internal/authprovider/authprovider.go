package authprovider

import (
	"github.com/H-BF/gw/internal/authprovider/casbinprovider"
	"github.com/H-BF/gw/pkg/authprovider"
)

func NewAuthProvider() (provider authprovider.AuthProvider, err error) {
	// todo: make auth plugin selection referring to conf

	provider, err = casbinprovider.NewCasbinAuthProvider("model.conf")

	return provider, err
}

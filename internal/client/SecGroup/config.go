package SecGroup

import "errors"

type AuthnType string

const (
	// AuthnTypeNONE - nonsecured - is by default
	AuthnTypeNONE AuthnType = "none"
	// AuthnTypeTLS - use TLS
	AuthnTypeTLS AuthnType = "tls"
)

var errAuthnTypeInvalid = errors.New("authn type invalid")

func (a *AuthnType) Validate() error {
	if !(*a == AuthnTypeNONE || *a == AuthnTypeTLS) {
		return errAuthnTypeInvalid
	}
	return nil
}

// TODO: вынести весь блок переменных в конфиг
var (
	sgroupsAddr = "localhost:9000"

	authnType         = AuthnTypeTLS
	authnKeyFile      = "/Users/yan/code/cert-workflow/client.key"
	authnCertFile     = "/Users/yan/code/cert-workflow/client.crt"
	authnServerVerify = true
	authnServerName   = "localhost"
	authnServerCAs    = []string{"/Users/yan/code/cert-workflow/chain.crt"}
)

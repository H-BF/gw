package gwconfig

import "github.com/H-BF/gw/internal/config"

/* config-sample.yaml

server:
	tls:
		disabled: true
		key-file: "/path/to/key-file.pem"
		cert-file: "/path/to/cert-file.pem"

external-api:
	sgroups:
		tls:
			disabled: true
			key-file: "/path/to/key-file.pem"
			cert-file: "/path/to/cert-file.pem"
*/

const (
	ServerTLSDisabled config.ValueT[bool]   = "server/tls/disabled"
	ServerTLSKeyFile  config.ValueT[string] = "server/tls/key-file"
	ServerTLSCertFile config.ValueT[string] = "server/tls/cert-file"

	ExternalApiSgroupsTLSDisabled config.ValueT[bool]   = "external-api/sgroups/tls/disabled"
	ExternalApiSgroupsTLSKeyFile  config.ValueT[string] = "external-api/sgroups/tls/key-file"
	ExternalApiSgroupsTLSCertFile config.ValueT[string] = "external-api/sgroups/tls/cert-file"
)

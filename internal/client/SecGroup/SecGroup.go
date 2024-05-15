package SecGroup

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/http2"
	"net"
	"net/http"

	gwconfig "github.com/H-BF/gw/internal/gw-config"

	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
)

type secGroupClient struct {
	sgroupsconnect.SecGroupServiceClient
}

func NewClient(addr string) sgroupsconnect.SecGroupServiceClient {
	transport := http.DefaultTransport

	if gwconfig.ExternalApiSgroupsTLSDisabled.MustValue(context.Background()) {
		transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		}
	} else {
		// TODO: enable tls for transport/http-client
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := sgroupsconnect.NewSecGroupServiceClient(
		httpClient,
		addr,
		connect.WithGRPC(),
	)

	return secGroupClient{client}
}

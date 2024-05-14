package SecGroup

import (
	"crypto/tls"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
	"golang.org/x/net/http2"
)

type secGroupClient struct {
	sgroupsconnect.SecGroupServiceClient
}

func NewClient(addr string) sgroupsconnect.SecGroupServiceClient {
	// client for transferring data via http2 to sgroups grpc server
	httpClient := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	client := sgroupsconnect.NewSecGroupServiceClient(
		httpClient,
		addr,
		connect.WithGRPC(),
	)

	return secGroupClient{client}
}

package SecGroup

import (
	"connectrpc.com/connect"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"os"
)

type secGroupClient struct {
	sgroupsconnect.SecGroupServiceClient
}

func NewClient() sgroupsconnect.SecGroupServiceClient {
	var addr string
	transport := &http2.Transport{}
	switch authnType {
	case AuthnTypeNONE:
		transport.AllowHTTP = true
		transport.DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		}
		addr = fmt.Sprintf("http://%s", sgroupsAddr)
	case AuthnTypeTLS:
		cfg, err := buildClientConf()
		if err != nil {
			panic(err)
		}
		transport.TLSClientConfig = cfg
		addr = fmt.Sprintf("https://%s", sgroupsAddr)
	}

	client := sgroupsconnect.NewSecGroupServiceClient(
		&http.Client{Transport: transport},
		addr,
		connect.WithGRPC(),
	)

	return secGroupClient{client}
}

func buildClientConf() (*tls.Config, error) {
	cfg := &tls.Config{}
	if authnKeyFile != "" && authnCertFile != "" {
		keyPair, err := tls.LoadX509KeyPair(authnCertFile, authnKeyFile)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = append(cfg.Certificates, keyPair)
	}
	cfg.InsecureSkipVerify = !authnServerVerify
	if authnServerVerify {
		certPool := x509.NewCertPool()
		for _, ca := range authnServerCAs {
			caData, err := os.ReadFile(ca)
			if err != nil {
				return nil, err
			}
			certPool.AppendCertsFromPEM(caData)
		}

		cfg.RootCAs = certPool
		cfg.ServerName = authnServerName
	}
	return cfg, nil
}

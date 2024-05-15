package httpserver

import (
	"context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"time"

	api "github.com/H-BF/gw/internal/api/SecGroup"
	"github.com/H-BF/gw/internal/authprovider"
	gwconfig "github.com/H-BF/gw/internal/gw-config"

	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
)

func ListenAndServe(addr string) error {
	casbinAuthProvider, err := authprovider.NewCasbinAuthProvider("model.conf", "policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(sgroupsconnect.NewSecGroupServiceHandler(
		api.NewSecGroupService(casbinAuthProvider),
	))

	var handler http.Handler
	if gwconfig.ServerTLSDisabled.MustValue(context.Background()) {
		handler = h2c.NewHandler(mux, &http2.Server{})
	} else {
		handler = mux
		// TODO: enable tls for server
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("server listening on http://127.0.0.1:%s", addr)

	return srv.ListenAndServe()
}

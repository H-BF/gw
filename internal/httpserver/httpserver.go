package httpserver

import (
	"log"
	"net/http"
	"time"

	api "github.com/H-BF/gw/internal/api/SecGroup"
	"github.com/H-BF/gw/internal/authprovider"

	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
)

func ListenAndServe(addr string) error {
	casbinAuthProvider, err := authprovider.NewCasbinAuthProvider("model.conf")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(sgroupsconnect.NewSecGroupServiceHandler(
		api.NewSecGroupService(casbinAuthProvider),
	))

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("server listening on http://127.0.0.1:%s", addr)

	return srv.ListenAndServe()
}

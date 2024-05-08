package main

import (
	"github.com/H-BF/protos/pkg/api/sgroups/sgroupsconnect"
	api "github.com/H-BF/sgroups-k8s-adapter/internal/api/SecGroup"
	"github.com/H-BF/sgroups-k8s-adapter/internal/authprovider"
	"log"
	"net/http"
)

func main() {
	casbinAuthProvider, err := authprovider.NewCasbinAuthProvider("model.conf", "policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle(sgroupsconnect.NewSecGroupServiceHandler(
		api.NewSecGroupService(casbinAuthProvider),
	))

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalln(err)
	}
}

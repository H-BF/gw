package main

import (
	"github.com/H-BF/gw/internal/config"
	gwconfig "github.com/H-BF/gw/internal/gw-config"
	"github.com/H-BF/gw/internal/httpserver"
	"log"
)

func main() {
	err := config.InitGlobalConfig(
		config.WithAcceptEnvironment{EnvPrefix: "GW"},
		config.WithSourceFile{FileName: ConfigFile},
		config.WithDefValue{Key: gwconfig.ServerTLSDisabled, Val: true},
		config.WithDefValue{Key: gwconfig.ExternalApiSgroupsTLSDisabled, Val: true},
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := httpserver.ListenAndServe(":8080"); err != nil {
		log.Fatalln(err)
	}
}

package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

var srv *http.Server

func ListenAndServe(port uint, handlers http.Handler) error {
	srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handlers,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("server listening on http://127.0.0.1:%d", port)

	return srv.ListenAndServe()
}

func Shutdown(ctx context.Context) error {
	return srv.Shutdown(ctx)
}

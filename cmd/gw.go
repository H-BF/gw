package main

import (
	"log"

	"github.com/H-BF/gw/internal/httpserver"
)

func main() {
	if err := httpserver.ListenAndServe(":8080"); err != nil {
		log.Fatalln(err)
	}
}

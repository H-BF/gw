package main

import (
	"github.com/H-BF/gw/internal/httpserver"
	"log"
)

func main() {
	if err := httpserver.ListenAndServe(":8080"); err != nil {
		log.Fatalln(err)
	}
}

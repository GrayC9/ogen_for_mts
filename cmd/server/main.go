//go:build ogen

package main

import (
	"log"
	"net/http"
	"os"

	api "github.com/example/ogen_for_mts/internal/api"
	"github.com/example/ogen_for_mts/internal/server"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}

	h, err := server.NewInMemoryHandler()
	if err != nil {
		log.Fatal(err)
	}

	s, err := api.NewServer(h)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, s); err != nil {
		log.Fatal(err)
	}
}

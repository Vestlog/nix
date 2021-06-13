package main

import (
	"log"
	"net/http"

	"github.com/vestlog/nix/pkg/api"
)

var (
	dsn = "storage.db"
)

func main() {
	handler, err := api.CreateAPIHandler(dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

package main

import (
	"fmt"
	"github.com/nvsoft/goapp/config"
	"log"
	"net/http"
)

func main() {
	err := startServer()

	if err != nil {
		log.Fatalf("Could not start HTTP server. %s", err)
	}
}

func startServer() error {
	staticHandler := http.FileServer(http.Dir(config.Config["static_root"]))
	http.Handle("/", http.StripPrefix("/", staticHandler))

	fmt.Printf("Starting server on port %s\n", config.Config["listen_address"])

	return http.ListenAndServe(config.Config["listen_address"], nil)
}

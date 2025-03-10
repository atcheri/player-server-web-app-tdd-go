package main

import (
	"log"
	"net/http"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

func main() {
	srv := &server.PlayerServer{}
	log.Fatal(http.ListenAndServe(":5000", srv))
}

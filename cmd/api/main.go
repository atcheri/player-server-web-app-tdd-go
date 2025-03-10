package main

import (
	"log"
	"net/http"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	memorystore "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

func main() {
	srv := &server.PlayerServer{Store: &memorystore.InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", srv))
}

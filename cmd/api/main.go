package main

import (
	"log"
	"net/http"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := persistence.LoadFileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer close()

	server := server.NewPlayerServer(store)
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}

}

package main

import (
	"log"
	"net/http"
	"os"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := &persistence.FileSystemPlayerStore{Database: db}
	srv := server.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":5000", srv))
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	poker "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/cli"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := persistence.LoadFileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	defer close()

	game := poker.NewCLI(store, os.Stdin, domain.BlindAlerterFunc(domain.StdOutAlerter))
	game.PlayPoker()
}

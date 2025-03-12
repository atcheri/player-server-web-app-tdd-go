package poker

import "github.com/atcheri/player-server-web-app-tdd-go/internal/domain"

type CLI struct {
	PlayerStore domain.PlayerStore
}

func (cli *CLI) PlayPoker() {
	cli.PlayerStore.RecordWin("Tester")
}

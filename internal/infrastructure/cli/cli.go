package poker

import (
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type CLI struct {
	PlayerStore domain.PlayerStore
	In          io.Reader
}

func (cli *CLI) PlayPoker() {
	cli.PlayerStore.RecordWin("Chris")
}

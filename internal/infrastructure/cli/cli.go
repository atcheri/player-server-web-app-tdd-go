package poker

import (
	"bufio"
	"io"
	"strings"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type CLI struct {
	PlayerStore domain.PlayerStore
	In          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.In)
	reader.Scan()

	cli.PlayerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

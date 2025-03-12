package poker

import (
	"bufio"
	"io"
	"strings"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type CLI struct {
	playerStore domain.PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store domain.PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}

func (cli *CLI) PlayPoker() {
	input := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(input))
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

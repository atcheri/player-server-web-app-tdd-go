package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

const PlayerPrompt = "Please enter the number of players: "

type CLI struct {
	playerStore domain.PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	game        *domain.Game
}

func NewCLI(store domain.PlayerStore, in io.Reader, out io.Writer, alerter domain.BlindAlerter) *CLI {
	game := &domain.Game{
		Alerter: alerter,
		Store:   store,
	}
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out:         out,
		game:        game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	players, _ := strconv.Atoi(cli.readLine())
	cli.game.Start(players)
	cli.game.Finish(cli.readLine())
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

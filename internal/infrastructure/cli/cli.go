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
	in   *bufio.Scanner
	out  io.Writer
	game domain.Game
}

func NewCLI(in io.Reader, out io.Writer, game domain.Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)
	players, err := strconv.Atoi(cli.readLine())
	if err != nil {
		return
	}

	cli.game.Start(players)
	cli.game.Finish(cli.readLine())
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

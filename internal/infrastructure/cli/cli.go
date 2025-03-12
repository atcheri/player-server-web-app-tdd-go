package poker

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type CLI struct {
	playerStore domain.PlayerStore
	in          *bufio.Scanner
	alerter     domain.BlindAlerter
}

func NewCLI(store domain.PlayerStore, in io.Reader, alerter domain.BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		alerter:     alerter,
	}
}

func (cli *CLI) PlayPoker() {
	cli.alerter.ScheduleAlertAt(5*time.Second, 100)
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

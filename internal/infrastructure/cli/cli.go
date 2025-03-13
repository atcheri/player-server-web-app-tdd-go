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
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}

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

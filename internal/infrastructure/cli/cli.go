package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type CLI struct {
	playerStore domain.PlayerStore
	in          *bufio.Scanner
	out         io.Writer
	alerter     domain.BlindAlerter
}

func NewCLI(store domain.PlayerStore, in io.Reader, out io.Writer, alerter domain.BlindAlerter) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
		out:         out,
		alerter:     alerter,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, "Please enter the number of players: ")
	cli.scheduleBlindAlerts()
	input := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(input))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

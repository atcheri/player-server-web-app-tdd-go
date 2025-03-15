package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

const PlayerPrompt = "Please enter the number of players: "

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
	fmt.Fprint(cli.out, PlayerPrompt)
	players, _ := strconv.Atoi(cli.readLine())
	cli.scheduleBlindAlerts(players)
	input := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(input))
}

func (cli *CLI) scheduleBlindAlerts(nb int) {
	blindIncrement := time.Duration(5+nb) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

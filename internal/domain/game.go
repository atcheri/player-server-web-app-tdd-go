package domain

import (
	"strings"
	"time"
)

type Game struct {
	Alerter BlindAlerter
	Store   PlayerStore
}

func (p *Game) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.Alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *Game) Finish(winner string) {
	p.Store.RecordWin(extractWinner(winner))
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

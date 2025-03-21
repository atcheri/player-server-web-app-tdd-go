package domain

import (
	"io"
	"strings"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, alertDestination io.Writer)
	Finish(winner string)
}

type TexasHoldem struct {
	Alerter BlindAlerter
	Store   PlayerStore
}

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		Alerter: alerter,
		Store:   store,
	}
}

func (p *TexasHoldem) Start(numberOfPlayers int, alertDestination io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.Alerter.ScheduleAlertAt(blindTime, blind, alertDestination)
		blindTime = blindTime + blindIncrement
	}
}

func (p *TexasHoldem) Finish(winner string) {
	p.Store.RecordWin(extractWinner(winner))
}

func extractWinner(input string) string {
	return strings.Replace(input, " wins", "", 1)
}

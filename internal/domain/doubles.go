package domain

import (
	"bytes"
	"fmt"
	"time"
)

var (
	DummySpyAlerter  = &SpyBlindAlerter{}
	DummyPlayerStore = &StubPlayerStore{}
	DummyStdIn       = &bytes.Buffer{}
	DummyStdOut      = &bytes.Buffer{}
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.Scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.Alerts = append(s.Alerts, ScheduledAlert{duration, amount})
}

type GameSpy struct {
	StartCalled     bool
	FinishCalled    bool
	NumberOfPlayers int
	Winner          string
}

func (s *GameSpy) Start(nb int) {
	s.StartCalled = true
	s.NumberOfPlayers = nb
}

func (s *GameSpy) Finish(w string) {
	s.FinishCalled = true
	s.Winner = extractWinner(w)
}

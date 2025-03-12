package poker_test

import (
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	poker "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/cli"
	"github.com/stretchr/testify/assert"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []domain.Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() domain.League {
	return s.league
}

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &poker.CLI{playerStore}
	cli.PlayPoker()

	assert.Equal(t, 1, len(playerStore.winCalls))
}

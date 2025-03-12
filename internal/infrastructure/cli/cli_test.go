package poker_test

import (
	"strings"
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
	t.Run("records Chris wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		cli := &poker.CLI{playerStore, in}

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("records Cleo wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}
		cli := &poker.CLI{playerStore, in}

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Cleo")
	})
}

func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	assert.Equal(t, 1, len(store.winCalls))
	assert.Equal(t, winner, store.winCalls[0])
}

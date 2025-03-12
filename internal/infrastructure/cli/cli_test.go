package poker_test

import (
	"strings"
	"testing"
	"time"

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

var dummySpyAlerter = &SpyBlindAlerter{}

type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{duration, amount})
}

func TestCLI(t *testing.T) {
	t.Run("records Chris wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("records Cleo wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in, dummySpyAlerter)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		assert.Equal(t, 1, len(blindAlerter.alerts), "expected a blind alert to be scheduled")
	})
}

func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	assert.Equal(t, 1, len(store.winCalls))
	assert.Equal(t, winner, store.winCalls[0])
}

package poker_test

import (
	"fmt"
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

		cases := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {
				alert := blindAlerter.alerts[i]
				assert.LessOrEqual(t, i, len(blindAlerter.alerts))
				assert.Equal(t, c.expectedAmount, alert.amount)
				assert.Equal(t, alert.scheduledAt, c.expectedScheduleTime)
			})
		}
	})
}

func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	assert.Equal(t, 1, len(store.winCalls))
	assert.Equal(t, winner, store.winCalls[0])
}

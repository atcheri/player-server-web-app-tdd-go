package poker_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	poker "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/cli"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Run("records Chris wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("7\nChris wins\n")
		stubStore := &domain.StubPlayerStore{}
		game := domain.NewTexasHoldem(domain.DummySpyAlerter, stubStore)
		cli := poker.NewCLI(in, domain.DummyStdOut, game)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, stubStore, "Chris")
	})

	t.Run("records Cleo wins from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("7\nCleo wins\n")
		stubStore := &domain.StubPlayerStore{}
		game := domain.NewTexasHoldem(domain.DummySpyAlerter, stubStore)
		cli := poker.NewCLI(in, domain.DummyStdOut, game)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, stubStore, "Cleo")
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		// arrange
		stdout := &bytes.Buffer{}
		game := &domain.GameSpy{}
		cli := poker.NewCLI(domain.DummyStdIn, stdout, game)

		// act
		cli.PlayPoker()

		// assert
		assert.Equal(t, poker.PlayerPrompt, stdout.String())
	})

	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		// arrange
		blindAlerter := &domain.SpyBlindAlerter{}
		game := domain.NewTexasHoldem(blindAlerter, domain.DummyPlayerStore)

		// act
		game.Start(5)

		cases := []domain.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		// assert
		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		// arrange
		blindAlerter := &domain.SpyBlindAlerter{}
		game := domain.NewTexasHoldem(blindAlerter, domain.DummyPlayerStore)

		// act
		game.Start(7)

		cases := []domain.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		// assert
		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		// arrange
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &domain.GameSpy{}
		cli := poker.NewCLI(in, stdout, game)

		// act
		cli.PlayPoker()

		// assert
		if game.StartCalled {
			t.Errorf("game should not have started")
		}
	})
}

func assertPlayerWin(t testing.TB, store *domain.StubPlayerStore, winner string) {
	t.Helper()

	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, winner, store.WinCalls[0])
}

func checkSchedulingCases(cases []domain.ScheduledAlert, t *testing.T, blindAlerter *domain.SpyBlindAlerter) {
	for i, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			alert := blindAlerter.Alerts[i]
			assert.LessOrEqual(t, i, len(blindAlerter.Alerts))
			assert.Equal(t, c.Amount, alert.Amount)
			assert.Equal(t, alert.At, c.At, fmt.Sprintf("alert %d was not scheduled %v", i, blindAlerter.Alerts))
		})
	}
}

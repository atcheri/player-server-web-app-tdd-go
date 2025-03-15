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
		cli := poker.NewCLI(domain.DummyPlayerStore, in, domain.DummyStdOut, domain.DummySpyAlerter)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, domain.DummyPlayerStore, "Chris")
	})

	t.Run("records Cleo wind from user's input", func(t *testing.T) {
		// arrange
		in := strings.NewReader("7\nCleo wins\n")
		cli := poker.NewCLI(domain.DummyPlayerStore, in, domain.DummyStdOut, domain.DummySpyAlerter)

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, domain.DummyPlayerStore, "Cleo")
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		// arrange
		stdout := &bytes.Buffer{}
		cli := poker.NewCLI(domain.DummyPlayerStore, domain.DummyStdIn, stdout, domain.DummySpyAlerter)

		// act
		cli.PlayPoker()

		// assert
		assert.Equal(t, poker.PlayerPrompt, stdout.String())
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &domain.SpyBlindAlerter{}

		cli := poker.NewCLI(domain.DummyPlayerStore, in, stdout, blindAlerter)
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []domain.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		for i, c := range cases {
			t.Run(fmt.Sprint(c), func(t *testing.T) {
				alert := blindAlerter.Alerts[i]
				assert.LessOrEqual(t, i, len(blindAlerter.Alerts))
				assert.Equal(t, c.Amount, alert.Amount)
				assert.Equal(t, alert.At, c.At, fmt.Sprintf("alert %d was not scheduled %v", i, blindAlerter.Alerts))
			})
		}
	})
}

func assertPlayerWin(t testing.TB, store *domain.StubPlayerStore, winner string) {
	t.Helper()

	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, winner, store.WinCalls[0])
}

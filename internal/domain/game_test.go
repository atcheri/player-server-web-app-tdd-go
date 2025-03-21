package domain_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		// arrange
		game := domain.NewTexasHoldem(domain.DummySpyAlerter, domain.DummyPlayerStore)

		// act
		game.Start(5, io.Discard)

		// assert
		cases := []domain.ScheduledAlert{
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
			t.Run(fmt.Sprintf("%d scheduled for %v", c.Amount, c.At), func(t *testing.T) {
				alert := domain.DummySpyAlerter.Alerts[i]
				assert.LessOrEqual(t, i, len(domain.DummySpyAlerter.Alerts))
				assert.Equal(t, c.Amount, alert.Amount)
				assert.Equal(t, alert.At, c.At)
			})
		}
	})
}

func TestGame_Finish(t *testing.T) {
	// arrange
	store := &domain.StubPlayerStore{}
	game := domain.NewTexasHoldem(domain.DummySpyAlerter, store)
	winner := "Ruth"

	// act
	game.Finish(winner)

	// assert
	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, winner, store.WinCalls[0])
}

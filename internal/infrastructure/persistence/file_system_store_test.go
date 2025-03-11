package persistence_test

import (
	"strings"
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
	"github.com/stretchr/testify/assert"

	filestore "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("get league from a file reader", func(t *testing.T) {
		// arrange
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`,
		)
		expectedLeague := []player.Player{
			{Name: "Cleo", Wins: 10},
			{Name: "Chris", Wins: 33},
		}
		store := filestore.FileSystemPlayerStore{database}

		// act
		league := store.GetLeague()

		// assert
		assert.Equal(t, expectedLeague, league)
		// reading again the league
		assert.Equal(t, expectedLeague, store.GetLeague())
	})
}

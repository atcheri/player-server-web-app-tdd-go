package persistence_test

import (
	"strings"
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
	"github.com/stretchr/testify/assert"
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
		league := store.getLeague()

		// assert
		assert.Equal(t, expectedLeague, league)

	})
}

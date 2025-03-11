package persistence_test

import (
	"io"
	"os"
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
	"github.com/stretchr/testify/assert"

	filestore "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("get league from a file reader", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`,
		)
		defer cleanDatabase()
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

	t.Run("get player score", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`,
		)
		defer cleanDatabase()
		store := filestore.FileSystemPlayerStore{database}

		// act
		score := store.GetPlayerScore("Chris")

		// assert
		assert.Equal(t, 33, score)

	})
}

func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

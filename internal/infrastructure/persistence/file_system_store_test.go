package persistence_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
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
		expectedLeague := domain.League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}
		store, _ := filestore.NewFileSystemPlayerStore(database)

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
		store, _ := filestore.NewFileSystemPlayerStore(database)

		// act
		score := store.GetPlayerScore("Chris")

		// assert
		assert.Equal(t, 33, score)
	})

	t.Run("store 2 wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, _ := filestore.NewFileSystemPlayerStore(database)

		store.RecordWin("Chris")
		store.RecordWin("Chris")

		assert.Equal(t, 35, store.GetPlayerScore("Chris"))
	})

	t.Run("store wins for new players", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store, _ := filestore.NewFileSystemPlayerStore(database)

		// act
		store.RecordWin("Pepper")

		// assert
		assert.Equal(t, store.GetPlayerScore("Pepper"), 1)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := filestore.NewFileSystemPlayerStore(database)

		assert.Nil(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := filestore.NewFileSystemPlayerStore(database)
		want := domain.League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		// act and assert
		assert.Nil(t, err)
		assert.Equal(t, want, store.GetLeague())
		assert.Equal(t, want, store.GetLeague())
	})
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
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

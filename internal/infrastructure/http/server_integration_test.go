package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	persistence "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
)

func TestRecordWinsAndRetrievePlayerScore(t *testing.T) {
	t.Run("get player Pepper's score", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		store, _ := persistence.NewFileSystemPlayerStore(database)
		server := server.NewPlayerServer(store)

		// act
		postRequest, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)

		response := httptest.NewRecorder()

		getRequest, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		server.ServeHTTP(response, getRequest)

		// assert
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), "3")
	})

	t.Run("get league players and their respective scores", func(t *testing.T) {
		// arrange
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		store, _ := persistence.NewFileSystemPlayerStore(database)
		server := server.NewPlayerServer(store)
		postRequest, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)

		// act
		request, _ := http.NewRequest(http.MethodPost, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		var playerJson []domain.Player
		err := json.NewDecoder(response.Body).Decode(&playerJson)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, []domain.Player{{Name: "Pepper", Wins: 3}}, playerJson)
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

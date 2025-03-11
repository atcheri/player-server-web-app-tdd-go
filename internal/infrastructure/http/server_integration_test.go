package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	memorystore "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
)

func TestRecordWinsAndRetrievePlayerScore(t *testing.T) {
	t.Run("get player Pepper's score", func(t *testing.T) {
		// arrange
		store := memorystore.NewInMemoryPlayerStore()
		server := server.NewPlayerServer(&store)

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
		store := memorystore.NewInMemoryPlayerStore()
		server := server.NewPlayerServer(&store)
		postRequest, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)
		server.ServeHTTP(httptest.NewRecorder(), postRequest)

		// act
		request, _ := http.NewRequest(http.MethodPost, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		var playerJson []player.Player
		err := json.NewDecoder(response.Body).Decode(&playerJson)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, []player.Player{{Name: "Pepper", Wins: 3}}, playerJson)
	})
}

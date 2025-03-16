package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

func TestGETPlayer(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		// arrange
		store := domain.StubPlayerStore{
			Scores: map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			WinCalls: nil, League: nil,
		}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		score := response.Body.String()
		expectedScore := "20"
		status := response.Code

		// assert
		assert.Equal(t, expectedScore, score)
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		// arrange
		store := domain.StubPlayerStore{
			Scores: map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			WinCalls: nil, League: nil,
		}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		score := response.Body.String()
		expectedScore := "10"
		status := response.Code

		// assert
		assert.Equal(t, expectedScore, score)
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("returns 404 no missing player", func(t *testing.T) {
		// arrange
		store := domain.StubPlayerStore{
			Scores: map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			WinCalls: nil, League: nil,
		}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodGet, "/players/NoOne", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		status := response.Code

		// assert
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestStorePlayerWins(t *testing.T) {
	t.Run("records win on POST request", func(*testing.T) {
		// arrange
		store := domain.StubPlayerStore{
			Scores:   map[string]int{},
			WinCalls: nil, League: nil,
		}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		status := response.Code

		// assert
		assert.Equal(t, http.StatusAccepted, status)
		AssertPlayerWins(t, store, "Pepper")
	})
}

func TestLeague(t *testing.T) {
	t.Run("returns 200 on /league", func(t *testing.T) {
		// arrange
		players := []domain.Player{
			{Name: "Cleo", Wins: 32},
			{Name: "Chris", Wins: 20},
			{Name: "Tiest", Wins: 14},
		}

		store := domain.StubPlayerStore{Scores: nil, WinCalls: nil, League: players}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodPost, "/league", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		status := response.Code
		var playerJson []domain.Player
		err := json.NewDecoder(response.Body).Decode(&playerJson)

		// assert
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "application/json", response.Result().Header.Get("content-type"))
		assert.Nil(t, err)
		assert.Equal(t, players, playerJson)
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		// arrange
		server := server.NewPlayerServer(&domain.StubPlayerStore{})
		request := newGameRequest()
		response := httptest.NewRecorder()

		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestWebSocket(t *testing.T) {
	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := domain.StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(server.NewPlayerServer(&store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
		}
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}

		AssertPlayerWins(t, store, winner)
	})
}

func AssertPlayerWins(t *testing.T, store domain.StubPlayerStore, winner string) {
	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, "Pepper", store.WinCalls[0])
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

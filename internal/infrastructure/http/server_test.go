package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
		srv, _ := server.NewPlayerServer(&store, &domain.GameSpy{})
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
		srv := mustMakePlayerServer(t, &store, &domain.GameSpy{})
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
		srv, _ := server.NewPlayerServer(&store, &domain.GameSpy{})
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
		srv := mustMakePlayerServer(t, &store, &domain.GameSpy{})
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
		srv := mustMakePlayerServer(t, &store, &domain.GameSpy{})
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
		server := mustMakePlayerServer(t, &domain.StubPlayerStore{}, &domain.GameSpy{})
		request := newGameRequest()
		response := httptest.NewRecorder()

		// act
		server.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestWebSocket(t *testing.T) {
	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		// arrange
		timeout := 10 * time.Millisecond
		expectedAlert := "Blind is 100"
		game := &domain.GameSpy{BlindAlert: []byte(expectedAlert)}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, domain.DummyPlayerStore, game))
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		// act
		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		// assert
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
		within(t, timeout, func() { assertWebsocketGotMsg(t, ws, expectedAlert) })
	})
}

func AssertPlayerWins(t *testing.T, store domain.StubPlayerStore, winner string) {
	assert.Equal(t, 1, len(store.WinCalls))
	assert.Equal(t, winner, store.WinCalls[0])
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func mustMakePlayerServer(t *testing.T, store domain.PlayerStore, game domain.Game) *server.PlayerServer {
	server, err := server.NewPlayerServer(store, game)
	assert.Nil(t, err, "problem creating player server")

	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.Nil(t, err, "could not open a ws connection on %s %v")

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func assertGameStartedWith(t *testing.T, game *domain.GameSpy, numberOfPlayers int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.NumberOfPlayers == numberOfPlayers
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", numberOfPlayers, game.NumberOfPlayers)
	}
}

func assertFinishCalledWith(t *testing.T, game *domain.GameSpy, winner string) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.Winner == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.Winner)
	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, expectedAlert string) {
	_, gotBlindAlert, _ := ws.ReadMessage()
	assert.Equal(t, expectedAlert, string(gotBlindAlert))
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

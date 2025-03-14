package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []domain.Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() domain.League {
	return s.league
}

func TestGETPlayer(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		// arrange
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			nil, nil,
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
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			nil, nil,
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
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			nil, nil,
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
		store := StubPlayerStore{
			map[string]int{},
			nil, nil,
		}
		srv := server.NewPlayerServer(&store)
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		status := response.Code

		// assert
		assert.Equal(t, http.StatusAccepted, status)
		assert.Equal(t, 1, len(store.winCalls))
		assert.Equal(t, "Pepper", store.winCalls[0])
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

		store := StubPlayerStore{nil, nil, players}
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

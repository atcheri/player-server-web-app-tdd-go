package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) error {
	s.winCalls = append(s.winCalls, name)
	return nil
}

func TestGETPlayer(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		// arrange
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
			nil,
		}
		srv := &server.PlayerServer{&store}
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
			nil,
		}
		srv := &server.PlayerServer{&store}
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
			nil,
		}
		srv := &server.PlayerServer{&store}
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
			nil,
		}
		srv := &server.PlayerServer{&store}
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		status := response.Code

		// assert
		assert.Equal(t, http.StatusAccepted, status)
		assert.Len(t, store.scores, 1)
		assert.Equal(t, "Pepper", store.winCalls[0])
	})
}

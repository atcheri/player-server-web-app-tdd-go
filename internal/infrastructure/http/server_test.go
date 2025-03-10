package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func TestGETPlayer(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		// arrange
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
		}
		srv := &server.PlayerServer{&store}
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		got := response.Body.String()
		want := "20"

		// assert
		assert.Equal(t, want, got)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		// arrange
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
		}
		srv := &server.PlayerServer{&store}
		request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		got := response.Body.String()
		want := "10"

		// assert
		assert.Equal(t, want, got)
	})

	t.Run("returns 404 no missing player", func(t *testing.T) {
		// arrange
		store := StubPlayerStore{
			map[string]int{
				"Pepper": 20,
				"Floyd":  10,
			},
		}
		srv := &server.PlayerServer{&store}
		request, _ := http.NewRequest(http.MethodGet, "/players/NoOne", nil)
		response := httptest.NewRecorder()

		// act
		srv.ServeHTTP(response, request)
		got := response.Code
		want := http.StatusNotFound

		// assert
		assert.Equal(t, want, got)
	})
}

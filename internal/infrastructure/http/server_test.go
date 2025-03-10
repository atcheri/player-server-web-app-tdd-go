package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
)

func TestGETPlayer(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		// arrange
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		// act
		server.PlayerServer(response, request)
		got := response.Body.String()
		want := "20"

		// assert
		assert.Equal(t, want, got)
	})
}

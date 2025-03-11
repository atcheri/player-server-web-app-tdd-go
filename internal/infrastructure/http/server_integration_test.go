package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	server "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/http"
	memorystore "github.com/atcheri/player-server-web-app-tdd-go/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
)

func TestRecordWinsAndRetrievePlayerScore(t *testing.T) {
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
}

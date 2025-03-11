package server

import (
	"fmt"
	"net/http"
	"strings"

	player "github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type PlayerServer struct {
	Store  player.PlayerStore
	router *http.ServeMux
}

func NewPlayerServer(store player.PlayerStore) *PlayerServer {
	server := &PlayerServer{
		store, http.NewServeMux(),
	}

	server.router.Handle("/league", http.HandlerFunc(server.handleLeague))
	server.router.Handle("/players/", http.HandlerFunc(server.handlePlayer))

	return server
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) handlePlayer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost:
		p.processPlayerWins(w, player)
	case http.MethodGet:
		p.showPlayerScore(w, player)
	}
}

func (p *PlayerServer) processPlayerWins(w http.ResponseWriter, player string) {
	_ = p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)

}

func (p *PlayerServer) showPlayerScore(w http.ResponseWriter, player string) {
	score := p.Store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

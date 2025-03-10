package server

import (
	"fmt"
	"net/http"
	"strings"

	player "github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type PlayerServer struct {
	Store player.PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.Store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

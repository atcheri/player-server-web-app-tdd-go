package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
	"github.com/gorilla/websocket"
)

const jsonContentType = "application/json"
const htmlTemplatePath = "game.html"

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type PlayerServer struct {
	Store domain.PlayerStore
	http.Handler
	template *template.Template
	Game     domain.Game
}

func NewPlayerServer(store domain.PlayerStore, game domain.Game) (*PlayerServer, error) {
	server := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	server.template = tmpl
	server.Store = store
	server.Game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(server.handleLeague))
	router.Handle("/players/", http.HandlerFunc(server.handlePlayer))
	router.Handle("/game", http.HandlerFunc(server.handleGame))
	router.Handle("/ws", http.HandlerFunc(server.handleWebSocket))

	server.Handler = router

	return server, nil
}

func (p *PlayerServer) handleLeague(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	leagueTable := p.getLeagueTable()
	json.NewEncoder(w).Encode(leagueTable)
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) getLeagueTable() domain.League {
	return p.Store.GetLeague()
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

func (p *PlayerServer) handleGame(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.Game.Start(numberOfPlayers, io.Discard) //todo: Don't discard the blinds messages!

	_, winnerMsg, _ := conn.ReadMessage()
	p.Game.Finish(string(winnerMsg))
}

func (p *PlayerServer) processPlayerWins(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)

}

func (p *PlayerServer) showPlayerScore(w http.ResponseWriter, player string) {
	score := p.Store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

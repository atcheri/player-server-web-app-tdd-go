package persistence

import (
	"sync"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type InMemoryPlayerStore struct {
	score map[string]int
	mu    sync.Mutex
}

func NewInMemoryPlayerStore() InMemoryPlayerStore {
	return InMemoryPlayerStore{score: make(map[string]int)}
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.score[name]
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.score[name]++
}

func (i *InMemoryPlayerStore) GetLeague() domain.League {
	var league []domain.Player

	for name, wins := range i.score {
		league = append(league, domain.Player{Name: name, Wins: wins})
	}

	return league
}

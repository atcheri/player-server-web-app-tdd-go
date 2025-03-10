package memorystore

import "sync"

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

func (s *InMemoryPlayerStore) RecordWin(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.score[name]++
	return nil
}

package memorystore

type InMemoryPlayerStore struct {
	score map[string]int
}

func NewInMemoryPlayerStore() InMemoryPlayerStore {
	return InMemoryPlayerStore{make(map[string]int)}
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return s.score[name]
}

func (s *InMemoryPlayerStore) RecordWin(name string) error {
	s.score[name]++
	return nil
}

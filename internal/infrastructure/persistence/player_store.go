package memorystore

type InMemoryPlayerStore struct{}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func (s *InMemoryPlayerStore) RecordWin(name string) error {
	return nil
}

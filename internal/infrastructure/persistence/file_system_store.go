package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type FileSystemPlayerStore struct {
	Database *json.Encoder
	league   domain.League
}

func NewFileSystemPlayerStore(database *os.File) (*FileSystemPlayerStore, error) {
	database.Seek(0, io.SeekStart)
	league, err := domain.NewLeague(database)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", database.Name(), err)
	}

	return &FileSystemPlayerStore{
		Database: json.NewEncoder(&Tape{database}),
		league:   league,
	}, nil
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	p := f.league.Find(name)
	if p != nil {
		p.Wins++
	} else {
		f.league = append(f.league, domain.Player{Name: name, Wins: 1})
	}

	f.Database.Encode(f.league)
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)
	if player == nil {
		return 0
	}

	return player.Wins
}

func (f *FileSystemPlayerStore) GetLeague() domain.League {
	return f.league
}

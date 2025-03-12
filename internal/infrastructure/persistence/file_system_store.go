package persistence

import (
	"encoding/json"
	"io"
	"os"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain"
)

type FileSystemPlayerStore struct {
	Database io.Writer
	league   domain.League
}

func NewFileSystemPlayerStore(database *os.File) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := domain.NewLeague(database)

	return &FileSystemPlayerStore{
		Database: &Tape{database},
		league:   league,
	}
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	p := f.league.Find(name)
	if p != nil {
		p.Wins++
	} else {
		f.league = append(f.league, domain.Player{Name: name, Wins: 1})
	}

	json.NewEncoder(f.Database).Encode(f.league)
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

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
	err := initialisePlayerDBFile(database)

	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

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

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}

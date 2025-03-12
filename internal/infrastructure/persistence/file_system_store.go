package persistence

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/league"
)

type FileSystemPlayerStore struct {
	Database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) RecordWin(name string) error {
	l := f.GetLeague()
	player := l.Find(name)
	if player == nil {
		return errors.New("player not found. Could not update score")
	}

	player.Wins++
	f.Database.Seek(0, io.SeekStart)
	json.NewEncoder(f.Database).Encode(l)
	return nil
}

func (f FileSystemPlayerStore) GetPlayerScore(name string) int {
	l := f.GetLeague()
	player := l.Find(name)
	if player == nil {
		return 0
	}

	return player.Wins
}

func (f *FileSystemPlayerStore) GetLeague() league.League {
	f.Database.Seek(0, io.SeekStart)
	league, _ := league.NewLeague(f.Database)
	return league
}

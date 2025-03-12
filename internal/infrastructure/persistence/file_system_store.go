package persistence

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/league"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type FileSystemPlayerStore struct {
	Database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) RecordWin(name string) error {
	league := f.GetLeague()
	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
			f.Database.Seek(0, io.SeekStart)
			json.NewEncoder(f.Database).Encode(league)
			return nil
		}
	}

	return errors.New("player not found. Could not update score")
}

func (f FileSystemPlayerStore) GetPlayerScore(name string) int {
	for _, player := range f.GetLeague() {
		if player.Name == name {
			return player.Wins
		}
	}

	return 0
}

func (f *FileSystemPlayerStore) GetLeague() []player.Player {
	f.Database.Seek(0, io.SeekStart)
	league, _ := league.NewLeague(f.Database)
	return league
}

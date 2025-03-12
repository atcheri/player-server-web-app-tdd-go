package persistence

import (
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/league"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type FileSystemPlayerStore struct {
	Database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) RecordWin(s string) error {
	return nil
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

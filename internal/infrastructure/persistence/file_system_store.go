package persistence

import (
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/league"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type FileSystemPlayerStore struct {
	Database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []player.Player {
	f.Database.Seek(0, io.SeekStart)
	league, _ := league.NewLeague(f.Database)
	return league
}

package persistence

import (
	"encoding/json"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type FileSystemPlayerStore struct {
	Database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []player.Player {
	var league []player.Player
	json.NewDecoder(f.Database).Decode(&league)
	return league
}

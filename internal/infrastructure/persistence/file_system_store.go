package persistence

import (
	"encoding/json"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/league"
	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type FileSystemPlayerStore struct {
	Database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) RecordWin(name string) {
	l := f.GetLeague()
	p := l.Find(name)
	if p != nil {
		p.Wins++
	} else {
		l = append(l, player.Player{Name: "Pepper", Wins: 1})
	}

	f.Database.Seek(0, io.SeekStart)
	json.NewEncoder(f.Database).Encode(l)
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

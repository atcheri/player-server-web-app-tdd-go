package league

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

type League []player.Player

func NewLeague(reader io.Reader) (League, error) {
	var league []player.Player
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}

func (l League) Find(name string) *player.Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}

	return nil
}

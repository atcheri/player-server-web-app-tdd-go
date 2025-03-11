package league

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/atcheri/player-server-web-app-tdd-go/internal/domain/player"
)

func NewLeague(reader io.Reader) ([]player.Player, error) {
	var league []player.Player
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}

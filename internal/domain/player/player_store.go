package player

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string) error
	GetLeague() []Player
}

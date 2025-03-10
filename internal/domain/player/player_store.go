package player

type PlayerStore interface {
	GetPlayerScore(name string) int
}

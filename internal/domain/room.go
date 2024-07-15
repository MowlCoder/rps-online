package domain

const (
	ROOM_WAITING_PLAYERS = iota
	ROOM_PLAYING
	ROOM_MATCH_END
)

type Room struct {
	ID       int
	Name     string
	Creator  User
	Opponent *User
	Status   uint8
	Match    *Match
}

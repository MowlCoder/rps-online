package domain

const (
	MATCH_WON = iota
	MATCH_LOST
	MATCH_DRAW
)

const (
	MATCH_BO1 = iota
)

type Match struct {
	Type          uint8
	CreatorCount  int
	OpponentCount int

	CurrentCreatorChoice  uint8
	CurrentOpponentChoice uint8

	LastCreatorChoice  uint8
	LastOpponentChoice uint8
}

func (m *Match) RoundsIsOver() bool {
	if m.Type == MATCH_BO1 {
		return m.CreatorCount+m.OpponentCount >= 1
	}

	return false
}

func (m *Match) ResetChoices() {
	m.LastCreatorChoice = m.CurrentCreatorChoice
	m.LastOpponentChoice = m.CurrentOpponentChoice

	m.CurrentCreatorChoice = CHOICE_NONE
	m.CurrentOpponentChoice = CHOICE_NONE
}

func (m *Match) BothPlayerHaveChosen() bool {
	return m.CurrentCreatorChoice != CHOICE_NONE && m.CurrentOpponentChoice != CHOICE_NONE
}

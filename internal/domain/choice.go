package domain

const (
	CHOICE_NONE = iota
	CHOICE_STONE
	CHOICE_PAPER
	CHOICE_SCISSORS
)

const (
	CHOICE_RESULT_DRAW = iota
	CHOICE_RESULT_LEFT
	CHOICE_RESULT_RIGHT
)

func ComputeChoiceResult(leftChoice uint8, rightChoice uint8) uint8 {
	if leftChoice == rightChoice {
		return CHOICE_RESULT_DRAW
	}

	if leftChoice == CHOICE_STONE {
		if rightChoice == CHOICE_PAPER {
			return CHOICE_RESULT_RIGHT
		}

		return CHOICE_RESULT_LEFT
	}

	if leftChoice == CHOICE_SCISSORS {
		if rightChoice == CHOICE_STONE {
			return CHOICE_RESULT_RIGHT
		}

		return CHOICE_RESULT_LEFT
	}

	if leftChoice == CHOICE_PAPER {
		if rightChoice == CHOICE_SCISSORS {
			return CHOICE_RESULT_RIGHT
		}

		return CHOICE_RESULT_LEFT
	}

	return CHOICE_RESULT_DRAW
}

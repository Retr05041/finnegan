package cellblock

type TimeLine struct {
	steps []CellBlock
}

type CellBlock struct {
	Candidate            string
	Direction            rune
	RowPlacedAt          int
	ColPlacedAt          int
	CandidatePlacedIndex int
}

func (t *TimeLine) NewStep(candidate string, direction rune, rowPlacedAt int, colPlacedAt int, candidatePlacedIndex int) {
	t.steps = append(t.steps, CellBlock{
		Candidate: candidate,
		Direction: direction,
		RowPlacedAt: rowPlacedAt,
		ColPlacedAt: colPlacedAt,
		CandidatePlacedIndex: candidatePlacedIndex,
	})
}

func (t *TimeLine) Backtrack() {
	t.steps = t.steps[:len(t.steps)-1]
}


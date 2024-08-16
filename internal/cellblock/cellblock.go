package cellblock

type CellBlock struct {
	Candidate            string
	Direction            rune
	RowPlacedAt          int
	ColPlacedAt          int
	CandidatePlacedIndex int
}

func NewCellBlock(candidate string, direction rune, rowPlacedAt int, colPlacedAt int, candidatePlacedIndex int) CellBlock {
	return CellBlock{
		Candidate: candidate,
		Direction: direction,
		RowPlacedAt: rowPlacedAt,
		ColPlacedAt: colPlacedAt,
		CandidatePlacedIndex: candidatePlacedIndex,
	}
}

package solver

import (
	"finnegan/internal/board"
	"fmt"
)

var (
	Timeline *board.Timeline
)

func Solve(t *board.Timeline) bool {
	Timeline = t
	if Solver() {
		return true
	}
	return false
}


// Main runner function
func Solver() bool {
	fmt.Println("Ready to start solving.")
	return false
}


// Idea: if I have the board keep track of it's "working cell" and "working direction", 
// each recursion step would do either a horizontal placement or a vertical placement
// keeping track of the candidates used (which should be "live" in reference to the global candidate list) 
// and where they were placed, each recursion step would create a "copy" of the board prior, so if we need to backtrack, 
// we simply pop the last board off the timeline and continue trying a new candidate like we never left the board in the first place

// The board should hold these values:
// grid :: current grid
// direction :: direction candidate was placed at
// WorkingRow :: working row cell for this board
// WorkingCol :: working col cell for this board 
// CandidatePlacementIndex :: index of candidate that goes in the working cell
// CellblockTotalLength :: length of the cell block
// CellblockBackLength :: 'left' / 'up' length coming off the working cell
// CellblockForwardLength :: 'right' / 'down' length coming off the working cell
// PossibleCandidates :: slice of possible candidates for this cell block :: used for backtracking (when a word here is used, it updated the timelines candidate map, if they become unused it adds it back to the candidatemap, this never gets candidated added back)
// Candidate :: Candidate used for this board

// The timeline should hold these values:
// CandidateMap :: Lengths of candidates -> slices of candidates :: This will be updated as candidates get used by a board
// CandidateReference :: perminant slice of candidates
// Boards :: slice of boards, the last board in the slice being the most recent

// In total we should have double the amount of boards in our timeline compared to white cells in the board
// Advancing brings us to the "next valid cell" which will be the first cell with a horizontal or vertical cellblock not filled / valid
// Backtracking brings us back to the board prior - the board holds possible candidates for that cellblock - inability to fit a valid candidate will result in a backtrack to happen again

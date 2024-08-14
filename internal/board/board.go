package board

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

// Current problems:
// 1. once all verticals are placed on a horizontal, there is no way to backtrack them...
// 2. No way to place any candidate that's not at the beginning index of that candidate

type Board struct {
	Size         int
	Grid         [][]rune
	CandidateMap map[int][]string
	DarkCell     rune
}

var (
	directions = [2]rune{'h', 'v'}
)

// Print the board nicely
func (b Board) Display() {
	for row := range b.Grid {
		for _, num := range b.Grid[row] {
			fmt.Print(string(num) + " ")
		}
		fmt.Println()
	}
}


func (b Board) SolveGame1Test() bool {
	// 70983
	_ = b.place("70983", 0, 2, 'h', 0)
	_ = b.place("701", 0, 2, 'v', 0)


	horizontalLengthOfWorkingCell, leftHorizontalLength, rightHorizontalLength := b.getHorizontalLength(1,3)
	if horizontalLengthOfWorkingCell == nil {
		return true
	}

	possibleHorizontalCandidates := b.CandidateMap[*horizontalLengthOfWorkingCell]
	for _, horizontalCandidate := range possibleHorizontalCandidates {
		fmt.Println("Current horizontal candidate: "+horizontalCandidate)
		if isHorizontalValid, _:= validHorizontalPlacement(b.Grid, horizontalCandidate, 1, 3, leftHorizontalLength, rightHorizontalLength); isHorizontalValid {
			_= b.place(horizontalCandidate, 1, 3, 'h', leftHorizontalLength)
			b.Display()
			fmt.Println("ready to move to the verticals!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
	return false
}

// Main runner function
func (b Board) Solve() bool {
	// STEPS:
	// 1. Go to next valid cell (non darkcell or border) - if the cell is the final cell then you have solved it
	// 2. Solve for horizontal (if it lands on a cell that's got a candidate in it skip it, otherwise it's not fully complete and an be filled) - if backtracked to here, use a different option
	// 3. Solve each vertical cell connected to that candidate - if any are invalid backtrack to step 2
	// 4. Repeat steps 1-4

	// Store this recursion steps used candidates
	var usedHorizontalCandidates []string
	var usedVerticalCandidates []string

	// Find the next cell for horizontal placement
	workingCellRow, workingCellCol := b.nextEmptyCell()
	if workingCellRow == nil || workingCellCol == nil {
		return true
	}
	fmt.Printf("WORKING HORZONTAL CELL: %d,%d\n", *workingCellRow, *workingCellCol)

	// Get the length of the current horizontal cell block
	horizontalLengthOfWorkingCell, leftHorizontalLength, rightHorizontalLength := b.getHorizontalLength(*workingCellRow, *workingCellCol)
	if horizontalLengthOfWorkingCell == nil {
		return true
	}

	// Loop through every candidate
	possibleHorizontalCandidates := b.CandidateMap[*horizontalLengthOfWorkingCell]
	for horizontalCanIndex, horizontalCandidate := range possibleHorizontalCandidates {
		// If it's already been tried skip
		if slices.Contains(usedHorizontalCandidates, horizontalCandidate) {
			continue
		}
		fmt.Println("Current horizontal candidate: "+horizontalCandidate)
		// If it's a Valid placement -- this will need to check back and forth, currently just checks forth
		if isHorizontalValid, unFilledCells := validHorizontalPlacement(b.Grid, horizontalCandidate, *workingCellRow, *workingCellCol, leftHorizontalLength, rightHorizontalLength); isHorizontalValid {
			usedHorizontalCandidates = append(usedHorizontalCandidates, horizontalCandidate)
			backupCells := b.place(horizontalCandidate, *workingCellRow, *workingCellCol, 'h', leftHorizontalLength)
			b.CandidateMap[*horizontalLengthOfWorkingCell] = removeCandidateFromList(possibleHorizontalCandidates, horizontalCanIndex)
			b.Display()
			fmt.Println("ready to move to the verticals!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			// While the verticals aren't solved.. -- validPlacement might need to return the locations of any characters that were freshly placed, as those are the ones in need of a vertical solve
			for _, cellColLocation := range unFilledCells {
				fmt.Printf("WORKING VERTICAL CELL: %d,%d\n", *workingCellRow, cellColLocation)

				verticalLengthOfWorkingCell := b.getPossibleVerticalLength(*workingCellRow, cellColLocation)
				if verticalLengthOfWorkingCell == nil {
					return true
				}

				possibleVerticalCandidates := b.CandidateMap[*verticalLengthOfWorkingCell]
				for verticalCanIndex, verticalCandidate := range possibleVerticalCandidates {
					if slices.Contains(usedVerticalCandidates, verticalCandidate) {
						continue
					}
					if validVerticalPlacement(b.Grid, verticalCandidate, *workingCellRow, cellColLocation) {
						usedVerticalCandidates = append(usedVerticalCandidates, verticalCandidate)
						backupVerticalCells := b.place(verticalCandidate, *workingCellRow, cellColLocation, 'v', 0)
						b.CandidateMap[*verticalLengthOfWorkingCell] = removeCandidateFromList(possibleVerticalCandidates, verticalCanIndex)
						b.Display()
						fmt.Println("Ready to move to the next empty cell!")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
						if b.Solve() {
							return true
						}
						b.remove(verticalCandidate, *workingCellRow, cellColLocation, 'v', backupVerticalCells, 0)
						b.CandidateMap[*verticalLengthOfWorkingCell] = addCandidateToList(possibleVerticalCandidates, verticalCandidate)
						fmt.Println("BACKTRACKED VERTICAL")
					}
				}
			}
			b.remove(horizontalCandidate, *workingCellRow, *workingCellCol, 'h', backupCells, leftHorizontalLength)
			b.CandidateMap[*horizontalLengthOfWorkingCell] = addCandidateToList(possibleHorizontalCandidates, horizontalCandidate)
			fmt.Println("BACKTRACKED HORIZONTAL")
		}
	}

	return false
}

// Checks if a candidate can be placed horizontally
func validHorizontalPlacement(grid [][]rune, candidate string, row int, col int, leftLength int, rightLength int) (bool, []int) {
	var unFilledCells []int
	if col < 0 || row < 0 || row >= len(grid) || col+rightLength > len(grid[row]) || col-leftLength < 0 {
		fmt.Printf("The candidate goes off the board - start col: %d, candidate length: %d, row length: %d\n", col, len(candidate), len(grid[row]))
		return false, nil
	}
	// If it's too small -- might be redundant
	if col+rightLength+1 < len(grid[row]) && grid[row][col+rightLength+1] == '.' || col-leftLength-1 > 0 && grid[row][col-leftLength-1] == '.' {
		fmt.Printf("It's too small to fit -- %d < %d && %t || %d > 0 && %t\n", col+rightLength, len(grid[row]), grid[row][col+rightLength+1] == '.', col-leftLength, grid[row][col-leftLength-1] == '.')
		return false, nil
	}
	// if it place nicely amongst the black cells and other words going right
	if rightLength > 0 {
		for r := range rightLength {
			nextCell := grid[row][col+r]
			if nextCell == '.' {
				unFilledCells = append(unFilledCells, col+r)
			}
			if nextCell != '.' && nextCell != rune(candidate[r]) {
				//fmt.Println("Next cell is not a '.' or the same as the current cell")
				return false, nil
			}
		}
	}
	// Going left
	if leftLength > 0 {
		for l := range leftLength {
			nextCell := grid[row][col-l]
			if nextCell == '.' {
				unFilledCells = append(unFilledCells, col-l)
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[l]) {
				//fmt.Println("Next cell is not a '.' or the same as the current cell")
				return false, nil
			}
		}
	}
	return true, unFilledCells
}

// Checks if a candidate can be placed vertically
func validVerticalPlacement(grid [][]rune, candidate string, row int, col int) bool {
	if col < 0 || row < 0 || col >= len(grid) || row+len(candidate)-1 > len(grid)-1 {
		fmt.Printf("The candidate goes off the board - start row: %d, candidate length: %d, col length: %d\n", row, len(candidate), col)
		return false
	}
	if row+len(candidate) < len(grid)-1 && grid[row+len(candidate)][col] == '.' {
		fmt.Println("It's too small to fit")
		return false
	}
	for i := range len(candidate) {
		nextCell := grid[row+i][col]
		if nextCell != '.' && nextCell != rune(candidate[i]) {
			//fmt.Println("Next cell is not a '.' or the same as the current cell")
			return false
		}
	}
	return true
}

// -- HELPER FUNCTIONS ---
func removeCandidateFromList(list []string, candidateIndex int) []string {
	list[candidateIndex] = list[len(list)-1]
	return list[:len(list)-1]
}

func addCandidateToList(list []string, candidate string) []string {
	return append(list, candidate)
}

// Get Horizontal length of cell we are on
func (b Board) getHorizontalLength(row int, col int) (*int,int,int) {
	totalLength := 1 // Account for current cell
	leftCell := 1
	rightCell := 1
	leftLength := 0
	rightLength := 0
	// Left
	for col-leftCell >= 0 && b.Grid[row][col-leftCell] != b.DarkCell {
		leftCell += 1
		leftLength += 1
	}
	// Right
	for col+rightCell < len(b.Grid) && b.Grid[row][col+rightCell] == '.' {
		rightCell += 1
		rightLength += 1
	}
	totalLength += leftLength + rightLength
	fmt.Printf("This cells allows for a horizontal size of %d - leftLength: %d - rightLength: %d\n", totalLength, leftLength, rightLength)
	if leftLength == 0 && rightLength == 0 {
		return nil, 0, 0
	}
	return &totalLength, leftLength, rightLength
}

// Get Vertical length of cell we are on
func (b Board) getPossibleVerticalLength(row int, col int) *int {
	totalLength := 1
	upCell := 1
	downCell := 1
	upLength := 0
	downLength := 0
	// Up
	for row-upCell >= 0 && b.Grid[row-upCell][col] != b.DarkCell {
		upCell += 1
		upLength += 1
	}
	// Down
	for row+downCell < len(b.Grid) && b.Grid[row+downCell][col] != b.DarkCell {
		downCell += 1
		downLength += 1
	}
	totalLength += downLength + upLength
	fmt.Printf("This cells allows for a vertical size of %d - upLength: %d - downLength: %d\n", totalLength, upLength, downLength)
	if upLength == 0 && downCell == 0 {
		return nil
	}
	return &totalLength
}

// Place a word horizontally or vertically at a specific row and col and return the updated Grid and a backup of what was in that cell block
func (b Board) place(candidate string, row int, col int, direction rune, startIndexOfCandidate int) []rune {
	var backupCellSequence []rune
	if direction == 'h' {
		for i := range len(candidate) {
			gridColOffset := col + (i - startIndexOfCandidate)
			if gridColOffset >= 0 && gridColOffset < len(b.Grid[row]) {
				backupCellSequence = append(backupCellSequence, b.Grid[row][gridColOffset])
				b.Grid[row][gridColOffset] = rune(candidate[i])
			}
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			gridRowOffset := row + (i - startIndexOfCandidate)
			if gridRowOffset >= 0 && gridRowOffset < len(b.Grid) && col >= 0 && col < len(b.Grid[gridRowOffset]) {
				backupCellSequence = append(backupCellSequence, b.Grid[gridRowOffset][col])
				b.Grid[gridRowOffset][col] = rune(candidate[i])
			}
		}
	}
	return backupCellSequence
}

// Remove a word, restore from backup, and return the updated Grid
func (b Board) remove(candidate string, row int, col int, direction rune, backupCellSequence []rune, startIndexOfCandidate int) {
	if direction == 'h' {
		for i := range len(candidate) {
			b.Grid[row][col+i] = backupCellSequence[i]
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			b.Grid[row+i][col] = backupCellSequence[i]
		}
	}
}

// Checks for next empty cell in the grid
func (b Board) nextEmptyCell() (*int, *int) {
	for row := range len(b.Grid) {
		for col := range len(b.Grid[row]) {
			if b.Grid[row][col] == '.' {
				return &row, &col
			}
		}
	}
	return nil, nil
}

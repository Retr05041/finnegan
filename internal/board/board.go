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

// Main runner function
func (b Board) Solve() bool {
	// STEPS:
	// 1. Go to next valid cell (non darkcell or border) - if the cell is the final cell then you have solved it
	// 2. Solve for horizontal (if it lands on a cell that's got a candidate in it skip it, otherwise it's not fully complete and an be filled) - if backtracked to here, use a different option
	// 3. Solve each vertical cell connected to that candidate - if any are invalid backtrack to step 2
	// 4. Repeat steps 1-4

	var usedCandidates []string
	workingCellRow, workingCellCol := b.nextEmptyCell()
	if workingCellRow == nil || workingCellCol == nil {
		return true
	}
	fmt.Printf("WORKING CELL H: %d,%d\n", *workingCellRow, *workingCellCol)

	horizontalLengthOfWorkingCell := b.getPossibleHorizontalLength(*workingCellRow, *workingCellCol)
	if horizontalLengthOfWorkingCell == nil {
		return true
	}

	possibleCandidates := b.CandidateMap[*horizontalLengthOfWorkingCell]
	for canIndex, candidate := range possibleCandidates {
		if slices.Contains(usedCandidates, candidate) {
			continue
		}
		if validPlacement(b.Grid, candidate, *workingCellRow, *workingCellCol, 'h') {
			usedCandidates = append(usedCandidates, candidate)
			newGrid, backupCells := place(b.Grid, candidate, *workingCellRow, *workingCellCol, 'h')
			b.Grid = newGrid
			b.CandidateMap[*horizontalLengthOfWorkingCell] = removeCandidateFromList(possibleCandidates, canIndex)
			b.Display()
			fmt.Println("Ready to move onto the next empty cell!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			if ! b.SolveVerticals(*workingCellRow, *workingCellCol) {
				b.Grid = remove(b.Grid, candidate, *workingCellRow, *workingCellCol, 'h', backupCells)
				b.CandidateMap[*horizontalLengthOfWorkingCell] = addCandidateToList(possibleCandidates, candidate)
				fmt.Println("BACKTRACKED")
				continue
			}
			if b.Solve() {
				return true
			}
			b.Grid = remove(b.Grid, candidate, *workingCellRow, *workingCellCol, 'h', backupCells)
			b.CandidateMap[*horizontalLengthOfWorkingCell] = addCandidateToList(possibleCandidates, candidate)
			fmt.Println("BACKTRACKED")
		}
	}

	return false
}

func (b Board) SolveVerticals(rowOfHorizontalCandidate int, colOfHorizontalCandidate int) bool {
	var usedCandidates []string
	workingCellRow, workingCellCol := rowOfHorizontalCandidate, colOfHorizontalCandidate
	if workingCellRow > len(b.Grid) || workingCellCol > len(b.Grid[workingCellRow]) {
		return true
	}
	if b.Grid[workingCellRow][workingCellCol] == b.DarkCell {
		return true
	}
	fmt.Printf("WORKING CELL V: %d,%d\n", workingCellRow, workingCellCol)
	
	verticalLengthOfWorkingCell := b.getPossibleVerticalLength(workingCellRow, workingCellCol)
	if verticalLengthOfWorkingCell == nil {
		return true
	}

	possibleCandidates := b.CandidateMap[*verticalLengthOfWorkingCell]
	for canIndex, candidate := range possibleCandidates {
		if slices.Contains(usedCandidates, candidate) {
			continue
		}
		fmt.Printf("Testing candidate: %s\n", candidate)
		if validPlacement(b.Grid, candidate, workingCellRow, workingCellCol, 'v') {
			usedCandidates = append(usedCandidates, candidate)
			newGrid, backupCells := place(b.Grid, candidate, workingCellRow, workingCellCol, 'v')
			b.Grid = newGrid
			b.CandidateMap[*verticalLengthOfWorkingCell] = removeCandidateFromList(possibleCandidates, canIndex)
			b.Display()
			fmt.Println("Ready to move onto the next vertical candidate")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			if b.SolveVerticals(workingCellRow, workingCellCol+1) {
				return true
			}
			b.Grid = remove(b.Grid, candidate, workingCellRow, workingCellCol, 'v', backupCells)
			b.CandidateMap[*verticalLengthOfWorkingCell] = addCandidateToList(possibleCandidates, candidate)
			fmt.Println("BACKTRACKED")
		}
	}
	return false
}

// Checks if a candidate can be placed at that location without breaking rules of the game - only works if the init row and col are empty
func validPlacement(grid [][]rune, candidate string, row int, col int, direction rune) bool {
	if direction == 'h' {
		if col < 0 || row < 0 || row >= len(grid) || col+len(candidate)-1 > len(grid[row]) {
			fmt.Printf("The candidate goes off the board - start col: %d, candidate length: %d, row length: %d\n", col, len(candidate), len(grid[row]))
			return false
		}
		// If it's too small
		if col+len(candidate) < len(grid[row]) && grid[row][col+len(candidate)] == '.' {
			fmt.Println("It's too small to fit")
			return false
		}
		// if it place nicely amongst the black cells and other words
		for i := range len(candidate) {
			nextCell := grid[row][col+i]
			if nextCell != '.' && nextCell != rune(candidate[i]) {
				fmt.Println("Next cell is not a '.' or the same as the current cell")
				return false
			}
			//if nextCell == rune(candidate[i]) && !canOverlapVertically(grid, candidate, row, col+i) {
			//	fmt.Println("Next cell is identical to what we want to place and it can't overlap the vertical word")
			//	return false
			//}
		}
	}

	if direction == 'v' {
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
				fmt.Println("Next cell is not a '.' or the same as the current cell")
				return false
			}
			//if nextCell == rune(candidate[i]) && !canOverlapHorizontally(grid, candidate, row+i, col) {
			//	fmt.Println("Next cell is identical to what we want to place and it can't overlap the vertical word")
			//	return false
			//}
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
func (b Board) getPossibleHorizontalLength(row int, col int) *int {
	totalLength := 1 // Account for current cell
	leftCell := 1
	rightCell := 1
	leftLength := 0
	rightLength := 0
	// Left
	for col-leftCell>= 0 && b.Grid[row][col-leftCell] != b.DarkCell {
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
		return nil
	}
	return &totalLength
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

// Checks if the candidate can be placed down on the cell block without breaking a pre-existing vertical word
func canOverlapVertically(grid [][]rune, candidate string, row int, col int) bool {
	for i := range len(candidate) {
		if grid[row+i][col] != '.' && grid[row+i][col] != rune(candidate[i]) {
			return false
		}
	}
	return true
}

// Checks if the candidate can be placed down on the cell block without breaking a pre-existing horizontal word
func canOverlapHorizontally(grid [][]rune, candidate string, row int, col int) bool {
	for i := range len(candidate) {
		if grid[row][col+i] != '.' && grid[row][col+i] != rune(candidate[i]) {
			return false
		}
	}
	return true
}

// Place a word horizontally or vertically at a specific row and col and return the updated Grid and a backup of what was in that cell block
func place(grid [][]rune, candidate string, row int, col int, direction rune) ([][]rune, []rune) {
	var backupCellSequence []rune
	if direction == 'h' {
		for i := range len(candidate) {
			backupCellSequence = append(backupCellSequence, grid[row][col+i])
			grid[row][col+i] = rune(candidate[i])
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			backupCellSequence = append(backupCellSequence, grid[row+i][col])
			grid[row+i][col] = rune(candidate[i])
		}
	}
	return grid, backupCellSequence
}

// Remove a word, restore from backup, and return the updated Grid
func remove(grid [][]rune, candidate string, row int, col int, direction rune, backupCellSequence []rune) [][]rune {
	if direction == 'h' {
		for i := range len(candidate) {
			grid[row][col+i] = backupCellSequence[i]
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			grid[row+i][col] = backupCellSequence[i]
		}
	}
	return grid
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

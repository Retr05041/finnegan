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
	// 701
	//_ = b.place("701", 0, 2, 'v', 0)

	fmt.Printf("WORKING VERTICAL CELL: %d,%d\n", 0, 2)

	verticalLengthOfWorkingCell, upLength, downLength := b.getPossibleVerticalLength(0, 2)
	if verticalLengthOfWorkingCell == nil {
		return true
	}

	possibleVerticalCandidates := b.CandidateMap[*verticalLengthOfWorkingCell]
	for verticalCanIndex, verticalCandidate := range possibleVerticalCandidates {
		fmt.Println("Current vertical candidate: " + verticalCandidate)
		if validVerticalPlacement(b.Grid, verticalCandidate, 0, 2, upLength, downLength) {
			backupVerticalCells := b.place(verticalCandidate, 0, 2, 'v', upLength)
			b.CandidateMap[*verticalLengthOfWorkingCell] = removeCandidateFromList(possibleVerticalCandidates, verticalCanIndex)
			b.Display()
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			b.remove(verticalCandidate, 0, 2, 'v', backupVerticalCells, upLength)
			b.CandidateMap[*verticalLengthOfWorkingCell] = addCandidateToList(possibleVerticalCandidates, verticalCandidate)
			fmt.Println("BACKTRACKED VERTICAL")
		}
	}

	return false
}

// Main runner function
func (b Board) Solve() bool {
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
		fmt.Println("Current horizontal candidate: " + horizontalCandidate)

		// If it's a Valid placement -- this will need to check back and forth, currently just checks forth
		if isHorizontalValid, unFilledCells := validHorizontalPlacement(b.Grid, horizontalCandidate, *workingCellRow, *workingCellCol, leftHorizontalLength, rightHorizontalLength); isHorizontalValid {
			usedHorizontalCandidates = append(usedHorizontalCandidates, horizontalCandidate)
			backupCells := b.place(horizontalCandidate, *workingCellRow, *workingCellCol, 'h', leftHorizontalLength)
			b.CandidateMap[*horizontalLengthOfWorkingCell] = removeCandidateFromList(possibleHorizontalCandidates, horizontalCanIndex)
			b.Display()
			fmt.Println("ready to move to the verticals!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			// ANY CELL THAT WAS JUST FILLED NEEDS TO BE VERTICIZED LOL (totally not listening to freebird while doing this bs)
			for _, cellColLocation := range unFilledCells {
				fmt.Printf("WORKING VERTICAL CELL: %d,%d\n", *workingCellRow, cellColLocation)

				verticalLengthOfWorkingCell, upLength, downLength := b.getPossibleVerticalLength(*workingCellRow, cellColLocation)
				if verticalLengthOfWorkingCell == nil {
					return true
				}

				possibleVerticalCandidates := b.CandidateMap[*verticalLengthOfWorkingCell]
				for verticalCanIndex, verticalCandidate := range possibleVerticalCandidates {
					if slices.Contains(usedVerticalCandidates, verticalCandidate) {
						continue
					}
					fmt.Println("Current vertical candidate: " + verticalCandidate)
					if validVerticalPlacement(b.Grid, verticalCandidate, *workingCellRow, cellColLocation, upLength, downLength) {
						usedVerticalCandidates = append(usedVerticalCandidates, verticalCandidate)
						backupVerticalCells := b.place(verticalCandidate, *workingCellRow, cellColLocation, 'v', upLength)
						b.CandidateMap[*verticalLengthOfWorkingCell] = removeCandidateFromList(possibleVerticalCandidates, verticalCanIndex)
						b.Display()
						fmt.Println("Ready to move to the next empty cell!")
						bufio.NewReader(os.Stdin).ReadBytes('\n')
						if b.Solve() {
							return true
						}
						b.remove(verticalCandidate, *workingCellRow, cellColLocation, 'v', backupVerticalCells, upLength)
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

	unFilledCells = append(unFilledCells, col) // Account for the cell we are currently on...

	// Left
	if leftLength > 0 {
		for l := 1; l <= leftLength; l++ {
			nextCell := grid[row][col-l]
			if nextCell == '.' {
				unFilledCells = append(unFilledCells, col-l)
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[l-1]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row, col-l)
				return false, nil
			}
		}
	}
	// Right
	if rightLength > 0 {
		for r := 1; r <= rightLength; r++ {
			nextCell := grid[row][col+r]
			if nextCell == '.' {
				unFilledCells = append(unFilledCells, col+r)
			}
			if nextCell != '.' && nextCell != rune(candidate[r-1]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row, col+r)
				return false, nil
			}
		}
	}
	return true, unFilledCells
}

// Checks if a candidate can be placed vertically
func validVerticalPlacement(grid [][]rune, candidate string, row int, col int, upLength int, downLength int) bool {
	if col < 0 || row < 0 || col >= len(grid) || row+downLength > len(grid)-1 || row-upLength < 0 {
		fmt.Printf("The candidate goes off the board - start row: %d, candidate length: %d, col length: %d\n", row, len(candidate), col)
		return false
	}
	if row+downLength+1 < len(grid)-1 && grid[row+downLength+1][col] == '.' || row-upLength-1 >= 0 && grid[row-upLength-1][col] == '.' {
		fmt.Println("It's too small to fit")
		return false
	}
	// Up
	if upLength > 0 {
		for u := 1; u <= upLength; u++ {
			nextCell := grid[row-u][col]
			if nextCell != '.' && nextCell != rune(candidate[u-1]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row-u, col)
				return false
			}
		}
	}
	// Down
	if downLength > 0 {
		for d := 1; d <= downLength; d++ {
			nextCell := grid[row+d][col]
			if nextCell != '.' && nextCell != rune(candidate[d-1]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row+d, col)
				return false
			}
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

func (b Board) getHorizontalLengths(row int, col int) (*int, int, int) {
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
	if leftLength == 0 && rightLength == 0 {
		return nil, 0, 0
	}
	return &totalLength, leftLength, rightLength
}

func (b Board) getVerticalLengths(row int, col int) (*int, int, int) {
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
		return nil, 0, 0
	}
	return &totalLength, upLength, downLength
}

func (b Board) placeHorizontal(candidate string, row int, col int, startIndexOfCandidate int) []rune {
	var backupCellSequence []rune
	for i := range len(candidate) {
		gridColOffset := col + (i - startIndexOfCandidate)
		if gridColOffset >= 0 && gridColOffset < len(b.Grid[row]) {
			backupCellSequence = append(backupCellSequence, b.Grid[row][gridColOffset])
			b.Grid[row][gridColOffset] = rune(candidate[i])
		}
	}
	return backupCellSequence
}

func (b Board) placeVertical(candidate string, row int, col int, startIndexOfCandidate int) []rune {
	var backupCellSequence []rune
	for i := range len(candidate) {
		gridColOffset := col + (i - startIndexOfCandidate)
		if gridColOffset >= 0 && gridColOffset < len(b.Grid[row]) {
			backupCellSequence = append(backupCellSequence, b.Grid[row][gridColOffset])
			b.Grid[row][gridColOffset] = rune(candidate[i])
		}
	}
	return backupCellSequence
}

func (b Board) removeHorizontal(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		b.Grid[row][col+i] = backupCellSequence[i]
	}
}


func (b Board) removeVertical(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		b.Grid[row+i][col] = backupCellSequence[i]
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

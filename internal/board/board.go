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
	_ = b.placeHorizontal("70983", 0, 2, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	backup1 := b.placeVertical("716", 0, 2, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	backup2 := b.placeVertical("013091", 0, 3, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	backup3 := b.placeHorizontal("110230", 1, 4, 2)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	b.removeHorizontal("110230", 1, 4, backup3, 2)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	b.removeVertical("013091", 0, 3, backup2, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	b.removeVertical("716", 0, 2, backup1, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')





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
	horizontalLengthOfCellBlock, leftLength, rightLength := b.getHorizontalLengths(*workingCellRow, *workingCellCol)
	if horizontalLengthOfCellBlock == nil {
		return true
	}
	fmt.Printf("This cell block as a horizontal length of: %d\n", *horizontalLengthOfCellBlock)

	// Loop through every candidate
	horizontalCandidates := b.CandidateMap[*horizontalLengthOfCellBlock]

	for horizontalCandidateIndex, horizontalCandidate := range horizontalCandidates {
		// If it's already been tried skip
		if slices.Contains(usedHorizontalCandidates, horizontalCandidate) {
			continue
		}

		fmt.Println("Current horizontal candidate: " + horizontalCandidate)

		// If it's a Valid placement
		if isHorizontalValid, _ := validHorizontalPlacement(b.Grid, horizontalCandidate, *workingCellRow, *workingCellCol, leftLength, rightLength); isHorizontalValid {
			usedHorizontalCandidates = append(usedHorizontalCandidates, horizontalCandidate)

			backupCells := b.placeHorizontal(horizontalCandidate, *workingCellRow, *workingCellCol, leftLength)
			b.CandidateMap[*horizontalLengthOfCellBlock] = removeCandidateFromList(horizontalCandidates, horizontalCandidateIndex)

			b.Display()
			fmt.Println("Ready to place the vertical at this cell!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			fmt.Printf("WORKING VERTICAL CELL: %d,%d\n", *workingCellRow, *workingCellCol)

			verticalLengthOfCellBlock, upLength, downLength := b.getVerticalLengths(*workingCellRow, *workingCellCol)
			if verticalLengthOfCellBlock == nil {
				return true
			}

			verticalCandidates := b.CandidateMap[*verticalLengthOfCellBlock]

			for verticalCandidateIndex, verticalCandidate := range verticalCandidates {
				if slices.Contains(usedVerticalCandidates, verticalCandidate) {
					continue
				}
				fmt.Println("Current vertical candidate: " + verticalCandidate)
				if validVerticalPlacement(b.Grid, verticalCandidate, *workingCellRow, *workingCellCol, upLength, downLength) {
					usedVerticalCandidates = append(usedVerticalCandidates, verticalCandidate)
					backupVerticalCells := b.placeVertical(verticalCandidate, *workingCellRow, *workingCellCol, upLength)
					b.CandidateMap[*verticalLengthOfCellBlock] = removeCandidateFromList(verticalCandidates, verticalCandidateIndex)
					b.Display()
					fmt.Println("Ready to move to the next empty cell!")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					if b.Solve() {
						return true
					}
					b.removeVertical(verticalCandidate, *workingCellRow, *workingCellCol, backupVerticalCells, upLength)
					b.CandidateMap[*verticalLengthOfCellBlock] = addCandidateToList(verticalCandidates, verticalCandidate)
					fmt.Println("BACKTRACKED VERTICAL")
				}
			}
			b.removeHorizontal(horizontalCandidate, *workingCellRow, *workingCellCol, backupCells, leftLength)
			b.CandidateMap[*horizontalLengthOfCellBlock] = addCandidateToList(horizontalCandidates, horizontalCandidate)
			fmt.Println("BACKTRACKED HORIZONTAL")
		}
	}

	return false
}

// Checks if a candidate can be placed horizontally
func validHorizontalPlacement(grid [][]rune, candidate string, row int, col int, leftLength int, rightLength int) (bool, []int) {
	var emptyCells []int

	if col < 0 || row < 0 || row >= len(grid) || col+rightLength > len(grid[row]) || col-leftLength < 0 {
		fmt.Printf("The candidate goes off the board - start col: %d, candidate length: %d, row length: %d\n", col, len(candidate), len(grid[row]))
		return false, nil
	}
	// If it's too small -- might be redundant
	if col+rightLength+1 < len(grid[row]) && grid[row][col+rightLength+1] == '.' || col-leftLength-1 > 0 && grid[row][col-leftLength-1] == '.' {
		fmt.Printf("It's too small to fit -- %d < %d && %t || %d > 0 && %t\n", col+rightLength, len(grid[row]), grid[row][col+rightLength+1] == '.', col-leftLength, grid[row][col-leftLength-1] == '.')
		return false, nil
	}

	emptyCells = append(emptyCells, col) // Account for the cell we are currently on...

	// Left
	if leftLength > 0 {
		for l := 1; l <= leftLength; l++ {
			nextCell := grid[row][col-l]
			if nextCell == '.' {
				emptyCells = append(emptyCells, col-l)
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
				emptyCells = append(emptyCells, col+r)
			}
			if nextCell != '.' && nextCell != rune(candidate[r-1]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row, col+r)
				return false, nil
			}
		}
	}
	return true, emptyCells
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
	// Check if the cell we are on is valid with what we want to place down
	if grid[row][col] != rune(candidate[upLength]) {
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
	for col+rightCell < len(b.Grid) && b.Grid[row][col+rightCell] != b.DarkCell {
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
		gridRowOffset := row + (i - startIndexOfCandidate)
		if gridRowOffset >= 0 && gridRowOffset < len(b.Grid)-1 {
			backupCellSequence = append(backupCellSequence, b.Grid[gridRowOffset][col])
			b.Grid[gridRowOffset][col] = rune(candidate[i])
		}
	}
	return backupCellSequence
}

func (b Board) removeHorizontal(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		gridColOffset := col + (i - startIndexOfCandidate)
		b.Grid[row][gridColOffset] = backupCellSequence[i]
	}
}


func (b Board) removeVertical(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		gridRowOffset := row + (i - startIndexOfCandidate)
		b.Grid[gridRowOffset][col] = backupCellSequence[i]
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

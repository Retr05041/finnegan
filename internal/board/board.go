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
	CurrentRow   int
	CurrentCol   int
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
	_ = b.placeVertical("716", 0, 2, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	_ = b.placeVertical("013091", 0, 3, 0)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	_ = b.placeHorizontal("110230", 1, 4, 2)
	b.Display()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	b.Display()
	_, leftLength, rightLength := b.getHorizontalLengths(2, 4)
	fmt.Printf("leftLength: %d - rightLength: %d", leftLength, rightLength)
	if b.validHorizontalPlacement("637047", 2, 4, leftLength, rightLength) {
		_ = b.placeHorizontal("637047", 2, 4, leftLength)
		b.Display()
	} else {
		fmt.Println("Can't fit..")
	}

	return false
}

// Main runner function
func (b Board) Solve() bool {
	// NEW STEPS
	// Will need a global timeline of type Cellblock that will hold all the info for a given h or v cellblock - this will be our backtracking format

	// next valid cell
	// horizontal check
		// if big enough to fit a word
			// if no open cells
				// Check if cellblock if valid
					// if not backtrack
			// fit word
				// if no words fit then backtrack
	// vertical check
		// identical to horizontal?

	
	nextValidRow, nextValidCol := b.nextValidCell()
	if nextValidRow == nil || nextValidCol == nil {
		return true
	}


	return false
}

// Checks if a candidate can be placed horizontally
func (b Board) validHorizontalPlacement(candidate string, row int, col int, leftLength int, rightLength int) bool {
	if col < 0 || row < 0 || row >= len(b.Grid) || col+rightLength > len(b.Grid[row]) || col-leftLength < 0 {
		fmt.Printf("The candidate goes off the board - start col: %d, candidate length: %d, row length: %d\n", col, len(candidate), len(b.Grid[row]))
		return false
	}
	// If it's too small -- might be redundant
	if col+rightLength+1 < len(b.Grid[row]) && b.Grid[row][col+rightLength+1] == '.' || col-leftLength-1 > 0 && b.Grid[row][col-leftLength-1] == '.' {
		fmt.Printf("It's too small to fit -- %d < %d && %t || %d > 0 && %t\n", col+rightLength, len(b.Grid[row]), b.Grid[row][col+rightLength+1] == '.', col-leftLength, b.Grid[row][col-leftLength-1] == '.')
		return false
	}

	// Left
	if leftLength > 0 {
		for l := 1; l <= leftLength; l++ {
			nextCell := b.Grid[row][col-l]
			if nextCell == '.' {
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[leftLength-l]) { // So close... this is broken...
				fmt.Printf("Cell %d,%d is not a '.' or %c\n", row, col-l, rune(candidate[leftLength-l]))
				return false
			}
		}
	}
	// Right
	if rightLength > 0 {
		for r := 1; r <= rightLength; r++ {
			nextCell := b.Grid[row][col+r]
			if nextCell == '.' {
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[leftLength+r]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row, col+r)
				return false
			}
		}
	}
	return true
}

// Checks if a candidate can be placed vertically
func (b Board) validVerticalPlacement(candidate string, row int, col int, upLength int, downLength int) bool {
	if col < 0 || row < 0 || col >= len(b.Grid) || row+downLength > len(b.Grid)-1 || row-upLength < 0 {
		fmt.Printf("The candidate goes off the board - start row: %d, candidate length: %d, col length: %d\n", row, len(candidate), col)
		return false
	}
	if row+downLength+1 < len(b.Grid)-1 && b.Grid[row+downLength+1][col] == '.' || row-upLength-1 >= 0 && b.Grid[row-upLength-1][col] == '.' {
		fmt.Println("It's too small to fit")
		return false
	}
	// Check if the cell we are on is valid with what we want to place down
	if b.Grid[row][col] != rune(candidate[upLength]) {
		return false
	}
	// Up
	if upLength > 0 {
		for u := 1; u <= upLength; u++ {
			nextCell := b.Grid[row-u][col]
			if nextCell == '.' {
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[upLength-u]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row-u, col)
				return false
			}
		}
	}
	// Down
	if downLength > 0 {
		for d := 1; d <= downLength; d++ {
			nextCell := b.Grid[row+d][col]
			if nextCell == '.' {
				continue
			}
			if nextCell != '.' && nextCell != rune(candidate[upLength+d]) {
				fmt.Printf("Cell %d,%d is not a '.' or the same as the current cell\n", row+d, col)
				return false
			}
		}
	}
	return true
}

// -- HELPER FUNCTIONS ---
func removeCandidateFromList(list []string, candidate string) []string {
	newList := []string{}
	for _, item := range list {
		if item != candidate {
			newList = append(newList, item)
		}
	}
	return newList
}

func addCandidateToList(list []string, candidate string) []string {
	for _, item := range list {
		if item == candidate {
			return list
		}
	}
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
	fmt.Printf("This cells allows for a total horizontal size of %d : leftLength: %d - rightLength: %d\n", totalLength, leftLength, rightLength)
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
	fmt.Printf("This cells allows for a total vertical size of %d : upLength: %d - downLength: %d\n", totalLength, upLength, downLength)
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
		if gridRowOffset >= 0 && gridRowOffset < len(b.Grid) {
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

// Checks for next empty cell in the b.Grid
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

func (b Board) nextValidCell() (*int, *int) {
	for row := range len(b.Grid) {
		for col := range len(b.Grid[row]) {
			if b.Grid[row][col] != b.DarkCell {
				return &row, &col
			}
		}
	}
	return nil, nil
}

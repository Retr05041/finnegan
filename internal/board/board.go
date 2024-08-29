package board

import (
	"fmt"
	"slices"
)

var (
	directions = [2]rune{'h', 'v'}
)

type Timeline struct {
	Length             int
	CurrentBoard       Board
	CandidateMap       map[int][]string
	CandidateReference []string
	Boards             []Board
}

func (t *Timeline) Advance() {
	newBoard := t.CurrentBoard
	t.Boards = append(t.Boards, newBoard)
	t.CurrentBoard = newBoard
	t.CurrentBoard.NextValidCell()
}

type Board struct {
	Grid     [][]rune
	DarkCell rune

	WorkingRow                 *int
	WorkingCol                 *int
	WorkingDirection           rune
	IndexOfCandidateWhenPlaced int
	CellblockTotalLength       int
	CellblockBackLength        int
	CellblockForwardLength     int

	PossibleCandidates []string
	WorkingCandidate   string
}

// Print the board nicely
func (b Board) Display() {
	for row := range b.Grid {
		for _, num := range b.Grid[row] {
			fmt.Print(string(num) + " ")
		}
		fmt.Println()
	}
}

func (b Board) VerticalCellBlockIsEmpty(row, col int) bool {
	_, upLength, downLength := b.GetVerticalLengths(*b.WorkingRow, *b.WorkingCol)
	if b.Grid[row][col] == '.' {
		return true
	}
	// Up
	for u := range upLength {
		if b.Grid[row-u][col] == '.' {
			return true
		}
	}
	// Down
	for d := range downLength {
		if b.Grid[row+d][col] == '.' {
			return true
		}
	}
	fmt.Println("Vertical cell block is not empty")
	return false
}

func (b Board) VerticalCellBlockIsValid(candidateRef []string) bool {
	row := *b.WorkingRow
	col := *b.WorkingCol
	upSideOfCandidate := ""
	placementChar := ""
	downSideOfCandidate := ""

	_, upLength, downLength := b.GetVerticalLengths(row, col)
	placementChar = string(b.Grid[row][col])

	// Up
	if upLength > 0 {
		for u := 1; u < upLength+1; u++ {
			upSideOfCandidate = string(b.Grid[row-u][col]) + upSideOfCandidate
		}
	}
	// Down
	if downLength > 0 {
		for d := 1; d < downLength+1; d++ {
			downSideOfCandidate = downSideOfCandidate + string(b.Grid[row+d][col])
		}
	}

	cellBlockCandidate := upSideOfCandidate + placementChar + downSideOfCandidate

	if slices.Contains(candidateRef, cellBlockCandidate) {
		fmt.Println("Contains candidate: " + cellBlockCandidate)
		return true
	}
	fmt.Println("Invalid candidate: " + cellBlockCandidate)
	return false
}

func (b Board) HorizontalCellBlockIsEmpty(row, col int) bool {
	_, leftLength, rightLength := b.GetHorizontalLengths(row, col)
	if b.Grid[row][col] == '.' {
		return true
	}
	// Left
	for l := range leftLength {
		if b.Grid[row][col-l] == '.' {
			return true
		}
	}
	// Right
	for r := range rightLength {
		if b.Grid[row][col+r] == '.' {
			return true
		}
	}
	fmt.Println("Horizontal cell block is not empty")
	return false
}

func (b Board) HorizontalCellBlockIsValid(candidateRef []string) bool {
	row := *b.WorkingRow
	col := *b.WorkingCol
	leftSideOfCandidate := ""
	placementChar := ""
	rightSideOfCandidate := ""

	_, leftLength, rightLength := b.GetHorizontalLengths(row, col)
	placementChar = string(b.Grid[row][col])

	// Left
	if leftLength > 0 {
		for l := 1; l < leftLength+1; l++ {
			leftSideOfCandidate = string(b.Grid[row][col-l]) + leftSideOfCandidate
		}
	}
	// Right
	if rightLength > 0 {
		for r := 1; r < rightLength+1; r++ {
			rightSideOfCandidate = rightSideOfCandidate + string(b.Grid[row][col+r])
		}
	}

	cellBlockCandidate := leftSideOfCandidate + placementChar + rightSideOfCandidate

	if slices.Contains(candidateRef, cellBlockCandidate) {
		fmt.Println("Contains candidate: " + cellBlockCandidate)
		return true
	}
	fmt.Println("Invalid candidate: " + cellBlockCandidate)
	return false
}

// Checks if a candidate can be placed horizontally
func (b Board) ValidHorizontalPlacement(candidate string, row int, col int, leftLength int, rightLength int) bool {
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
func (b Board) ValidVerticalPlacement(candidate string, row int, col int, upLength int, downLength int) bool {
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
func RemoveCandidateFromList(list []string, candidate string) []string {
	newList := []string{}
	for _, item := range list {
		if item != candidate {
			newList = append(newList, item)
		}
	}
	return newList
}

func AddCandidateToList(list []string, candidate string) []string {
	for _, item := range list {
		if item == candidate {
			return list
		}
	}
	return append(list, candidate)
}

func (b Board) GetHorizontalLengths(row int, col int) (*int, int, int) {
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

func (b Board) GetVerticalLengths(row int, col int) (*int, int, int) {
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
	if upLength == 0 && downCell == 0 {
		return nil, 0, 0
	}
	return &totalLength, upLength, downLength
}

func (b Board) PlaceHorizontal(candidate string, row int, col int, startIndexOfCandidate int) []rune {
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

func (b Board) PlaceVertical(candidate string, row int, col int, startIndexOfCandidate int) []rune {
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

func (b Board) RemoveHorizontal(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		gridColOffset := col + (i - startIndexOfCandidate)
		b.Grid[row][gridColOffset] = backupCellSequence[i]
	}
}

func (b Board) RemoveVertical(candidate string, row int, col int, backupCellSequence []rune, startIndexOfCandidate int) {
	for i := range len(candidate) {
		gridRowOffset := row + (i - startIndexOfCandidate)
		b.Grid[gridRowOffset][col] = backupCellSequence[i]
	}
}

func (b *Board) NextValidCell() {
	for row := range len(b.Grid) {
		for col := range len(b.Grid[row]) {
			if b.Grid[row][col] != b.DarkCell && !b.HorizontalCellBlockIsEmpty(row, col) || !b.VerticalCellBlockIsEmpty(row, col) {
				b.WorkingRow = &row
				b.WorkingCol = &col
				return
			}
		}
	}
}

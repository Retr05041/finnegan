package board

import (
	"fmt"
	"bufio"
	"os"
)

// CURRENT PROBLEM - Candidates of the wrong size are being placed (smaller than whats needed) + if they backtrack they need to add the candidate back into the list... + I also believe it can't handle middle grid setting


type Board struct {
	Size int
	Grid [][]rune
	CandidateMap map[int][]string
	DarkCell rune
}

var (
	directions = [2]rune{'h','v'}
)

func (b Board) Display() {
	for row := range b.Grid {
		for _, num := range b.Grid[row] {
			fmt.Print(string(num) + " ")
		}
		fmt.Println()
	}
}

func (b Board) Solve() bool {
	emptyCellRow, emptyCellCol := nextEmptyCell(b.Grid)
	if emptyCellRow == nil || emptyCellCol == nil {
		return true
	}

	possibleCandidateLengths := getPossibleCandidateLengths(b.Grid, *emptyCellRow, *emptyCellCol)
	fmt.Printf("Possible candidates given the starting cell: %d,%d -- ", *emptyCellRow, *emptyCellCol)

	for _,length := range possibleCandidateLengths {
		possibleCandidates := b.CandidateMap[length]
		fmt.Println(possibleCandidates)
		fmt.Print("Untested lengths: ")
		fmt.Println(possibleCandidateLengths)
		for canIndex,candidate := range possibleCandidates {
			for _,direction := range directions {
				if validPlacement(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction) {
					newGrid, backupCells := place(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction)
					b.Grid = newGrid
					b.CandidateMap[length] = removeCandidateFromList(possibleCandidates, canIndex)
					b.Display()
					bufio.NewReader(os.Stdin).ReadBytes('\n') 
					if b.Solve() {
						return true
					}
					b.Display()
					b.Grid = remove(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction, backupCells)
					b.CandidateMap[length] = addCandidateToList(possibleCandidates, candidate)
				}
			}
		}
	}

	return false
}

func removeCandidateFromList(list []string, candidateIndex int) []string {
	list[candidateIndex] = list[len(list)-1]
	return list[:len(list)-1]
}

func addCandidateToList(list []string, candidate string) []string {
	return append(list, candidate)
}

// Checks for possible lengths given an empty cell
func getPossibleCandidateLengths(grid [][]rune, row int, col int) []int {
	var possibleLengths []int

	for _, direction := range directions {
		if direction == 'h' {
			if col > 0 && grid[row][col-1] != '.' {
				continue
			}
			length := 0
			for col + length < len(grid) && grid[row][col + length] == '.' {
				length += 1
			}
			if length > 0 {
				possibleLengths = append(possibleLengths, length)
			}
		}
		if direction == 'v' {
			if row > 0 && grid[row-1][col] != '.' {
				continue
			}
			length := 0
			for row + length < len(grid) && grid[row+length][col] == '.' {
				length += 1
			}
			if length > 0 {
				possibleLengths = append(possibleLengths, length)
			}
		}

	}
	return possibleLengths
}

// Checks if a candidate can be placed at that location without breaking rules of the game
func validPlacement(grid [][]rune, candidate string, row int, col int, direction rune) bool {
	if direction == 'h' {
		if col + len(candidate) > len(grid)-1 {
			return false
		}
		for i := range len(candidate) {
			cell := grid[row][col+i]
			// Already placed word / invalid placement check
			if cell != '.' && cell != rune(candidate[i]) {
				return false
			}
			if cell == rune(candidate[i]) && ! canOverlapVertically(grid, candidate, row, col+i) {
				return false
			}
		}
	}

	if direction == 'v' {
		if col + len(candidate) > len(grid)-1 {
			return false
		}
		for i := range len(candidate) {
			cell := grid[row+i][col]
			if cell != '.' && cell != rune(candidate[i]) {
				return false
			}
			if cell == rune(candidate[i]) && ! canOverlapHorizontally(grid, candidate, row+i, col) {
				return false
			}
		}
	}
	return true
}

func canOverlapVertically(grid [][]rune, candidate string, row int, col int) bool {
	for i := range len(candidate) {
		if grid[row+i][col] != '.' && grid[row+i][col] != rune(candidate[i]) {
			return false
		}
	}
	return true
}

func canOverlapHorizontally(grid [][]rune, candidate string, row int, col int) bool {
	for i := range len(candidate) {
		if grid[row][col+i] != '.' && grid[row][col+i] != rune(candidate[i]) {
			return false
		}
	}
	return true
}

// Place a word and return the updated Grid
func place(grid [][]rune, candidate string, row int, col int, direction rune) ([][]rune,[]rune){
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

// Remove a word and return the updated Grid
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

func nextEmptyCell(grid [][]rune) (*int,*int) {
	for row := range len(grid) {
		for col := range len(grid[row]) {
			if grid[row][col] == '.' {
				return &row, &col
			}
		}
	}
	return nil, nil
}


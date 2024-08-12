package board

import (
	"fmt"
)

// CURRENT PROBLEM - Candidates of the wrong size are being placed (smaller than whats needed) + if they backtrack they need to add the candidate back into the list... + I also believe it can't handle middle grid setting


type Board struct {
	Size int
	Grid [][]rune
	NumberList []string
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
	emptyCellRow, exmptyCellCol := nextEmptyCell(b.Grid)
	if emptyCellRow == nil || exmptyCellCol == nil {
		return true
	}

	for index, candidate := range b.NumberList {
		for _,direction := range directions {
			if validPlacement(b.Grid, candidate, *emptyCellRow, *exmptyCellCol, direction) {
				b.Grid = place(b.Grid, candidate, *emptyCellRow, *exmptyCellCol, direction)
				b.NumberList = useCandidate(b.NumberList,index)
				b.Display()
				if b.Solve() {
					return true
				}
				b.Grid = remove(b.Grid, candidate, *emptyCellRow, *exmptyCellCol, direction)
			}
		}
	}
	
	return false
}

func useCandidate(candidateList []string, candidateIndex int) []string {
	return append(candidateList[:candidateIndex], candidateList[candidateIndex+1:]...)
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
		}
		if grid[row][col+len(candidate)+1] == '.' {
			return false
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
		}
		if grid[row+len(candidate)+1][col] == '.' {
			return false
		}
	}
	return true
}

// Place a word and return the updated Grid
func place(grid [][]rune, candidate string, row int, col int, direction rune) [][]rune {
	if direction == 'h' {
		for i := range len(candidate) {
			grid[row][col+i] = rune(candidate[i])
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			grid[row+i][col] = rune(candidate[i])
		}
	}
	return grid
}

// Remove a word and return the updated Grid
func remove(grid [][]rune, candidate string, row int, col int, direction rune) [][]rune {
	if direction == 'h' {
		for i := range len(candidate) {
			grid[row][col+i] = '.'
		}
	} else if direction == 'v' {
		for i := range len(candidate) {
			grid[row+i][col] = '.'
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


package board

import (
	"bufio"
	"fmt"
	"os"
)

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
	emptyCellRow, emptyCellCol := nextEmptyCell(b.Grid)
	if emptyCellRow == nil || emptyCellCol == nil {
		return true
	}
	fmt.Printf("WORKING CELL: %d,%d", *emptyCellRow, *emptyCellCol)
	fmt.Println()

	possibleCandidateLengths := b.getPossibleCandidateLengths(*emptyCellRow, *emptyCellCol)
	fmt.Print("Possible lengths: ")
	fmt.Println(possibleCandidateLengths)

	for _, length := range possibleCandidateLengths {
		possibleCandidates := b.CandidateMap[length]
		fmt.Print("Testing candidates - ")
		fmt.Println(possibleCandidates)
		for canIndex, candidate := range possibleCandidates {
			fmt.Print("Testing candidate - ")
			fmt.Println(candidate)
			for _, direction := range directions {
				if validPlacement(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction) {
					newGrid, backupCells := place(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction)
					b.Grid = newGrid
					b.CandidateMap[length] = removeCandidateFromList(possibleCandidates, canIndex)
					b.Display()
					fmt.Println("Ready to move onto the next empty cell!")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					if b.Solve() {
						return true
					}
					b.Grid = remove(b.Grid, candidate, *emptyCellRow, *emptyCellCol, direction, backupCells)
					b.CandidateMap[length] = addCandidateToList(possibleCandidates, candidate)
					fmt.Println("BACKTRACKED")
				}
			}
		}
	}

	return false
}

// Checks if a candidate can be placed at that location without breaking rules of the game
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
				fmt.Println("Next cell is not a '.' or the same as the current cell?")
				return false
			}
			if nextCell == rune(candidate[i]) && !canOverlapVertically(grid, candidate, row, col+i) {
				fmt.Println("Next cell is identical to what we want to place and it can't overlap the vertical word")
				return false
			}
		}
	}

	if direction == 'v' {
		if row+len(candidate) > len(grid)-1 {
			return false
		}
		for i := range len(candidate) {
			nextCell := grid[row+i][col]
			if nextCell != '.' && nextCell != rune(candidate[i]) {
				return false
			}
			if nextCell == rune(candidate[i]) && !canOverlapHorizontally(grid, candidate, row+i, col) {
				return false
			}
		}
		if grid[row+len(candidate)][col] == '.' {
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

// Checks for possible lengths given an empty cell
func (b Board) getPossibleCandidateLengths(row int, col int) []int {
	var possibleLengths []int
	var length int

	for _, direction := range directions {
		if direction == 'h' {
			if col > 0 && b.Grid[row][col-1] != '.' && b.Grid[row][col-1] != b.DarkCell {
				continue
			}
			length = 0
			for col+length < len(b.Grid) && b.Grid[row][col+length] == '.' {
				length += 1
			}
			if length > 0 {
				fmt.Printf("This cells allows for a horizontal size of %d", length)
				fmt.Println()
				possibleLengths = append(possibleLengths, length)
			}
		}
		if direction == 'v' {
			length = 0
			for row-length < len(b.Grid) && b.Grid[row-length][col] != b.DarkCell && b.Grid[row-length][col] != '.' {
				length += 1
			}
			for row+length < len(b.Grid) && b.Grid[row+length][col] == '.' {
				length += 1
			}
			if length > 0 {
				fmt.Printf("This cells allows for a vertical size of %d", length)
				fmt.Println()
				possibleLengths = append(possibleLengths, length)
			}
		}

	}
	return possibleLengths
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
func nextEmptyCell(grid [][]rune) (*int, *int) {
	for row := range len(grid) {
		for col := range len(grid[row]) {
			if grid[row][col] == '.' {
				return &row, &col
			}
		}
	}
	return nil, nil
}

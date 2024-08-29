package board

import (
	"fmt"
	"os"
	"slices"
	"bufio"
)

// Main runner function
func (b *Board) Solve() bool {
	var triedCandidates []string

	b.nextFillableCell() // Move down the line (Will remain the same if there is a valid horizontal or vertical cellblock to be filled)

	fmt.Printf("Current working cell: %d,%d\n", *b.CurrentRow, *b.CurrentCol)

	if b.isHorizontalCellBlockIsEmpty(*b.CurrentRow, *b.CurrentCol) {
		fmt.Println("Horizontal has an empty cell in the block")
		totalHorizontalLength, leftLength, rightLength := b.getHorizontalLengths(*b.CurrentRow, *b.CurrentCol)
		if _, ok := b.CandidateMap[*totalHorizontalLength]; ok {
			for _, candidate := range b.CandidateMap[*totalHorizontalLength] { // Get a list of all possible candidates
				if slices.Contains(triedCandidates, candidate) { // Skip tried candidates
					continue
				}
				if b.isValidHorizontalPlacement(candidate, *b.CurrentRow, *b.CurrentCol, leftLength, rightLength) {
					horizontalBackup := b.placeHorizontalOnBoard(candidate, *b.CurrentRow, *b.CurrentCol, leftLength)
					triedCandidates = append(triedCandidates, candidate)
					b.Display()
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					if b.Solve() {
						return true
					}
					b.removeHorizontalFromBoard(candidate, *b.CurrentRow, *b.CurrentCol, horizontalBackup, leftLength)
				}
			}
		} else {
			fmt.Println("Horizontal cellblock is empty but is too short!")		
		}
	} else if !b.isHorizontalCellBlockIsValid(*b.CurrentRow, *b.CurrentCol) {
		fmt.Println("Horizontal cell block is not empty and does not contain a valid candidate")
		return false
	} else {
		if b.isVerticalCellBlockIsEmpty(*b.CurrentRow, *b.CurrentCol) {
			fmt.Println("Vertical has an empty cell in the block")
			totalVerticalLength, upLength, downLength := b.getVerticalLengths(*b.CurrentRow, *b.CurrentCol)
			for _, candidate := range b.CandidateMap[*totalVerticalLength] { // Get a list of all possible candidates
				if slices.Contains(triedCandidates, candidate) { // Skip tried candidates
					continue
				}
				if b.isValidVerticalPlacement(candidate, *b.CurrentRow, *b.CurrentCol, upLength, downLength) {
					verticalBackup := b.placeVerticalOnBoard(candidate, *b.CurrentRow, *b.CurrentCol, upLength)
					triedCandidates = append(triedCandidates, candidate)
					b.Display()
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					if b.Solve() {
						return true
					}
					b.removeVerticalFromBoard(candidate, *b.CurrentRow, *b.CurrentCol, verticalBackup, upLength)
				}
			}
		} else if !b.isVerticalCellBlockIsValid(*b.CurrentRow, *b.CurrentCol) {
			fmt.Println("Vertical cell block is not empty and does not contain a valid candidate")
			return false
		}
	}

	return false
}

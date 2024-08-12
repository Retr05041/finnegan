package reader

import (
	"bufio"
	"finnegan/internal/board"
	"os"
	"strconv"
)

func Read(gamePath string) (*board.Board, error) {
	// Open File
	file, err := os.Open(gamePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	newBoard := new(board.Board)
	newBoard.DarkCell = '■'
	newBoard.CandidateMap = make(map[int][]string)

	// Get Grid size
	gridSizeLine, err := fileReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	newBoard.Size, err = strconv.Atoi(gridSizeLine[:len(gridSizeLine)-1])
	if err != nil {
		return nil, err
	}

	// Initilize Grid
	newBoard.Grid = make([][]rune, newBoard.Size)
	for i := range newBoard.Grid {
		newBoard.Grid[i] = make([]rune, newBoard.Size)
	}

	// Gather grid from file
	for row := range newBoard.Size {
		gridLine, err := fileReader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		gridLine = gridLine[:len(gridLine)-1]
		runeLine := []rune(gridLine)
		for i := range runeLine {
			if runeLine[i] == '1' {
				runeLine[i] = newBoard.DarkCell
			} else {
				runeLine[i] = '.'
			}
		}
		newBoard.Grid[row] = runeLine
	}

	// Gather numbers from file
	for {
        line, err := fileReader.ReadString('\n')
        if err != nil {
            if err.Error() != "EOF" {
				return nil, err
            }
            break
        }
		line = line[:len(line)-1]

		if existingCandidates, ok := newBoard.CandidateMap[len(line)]; ok {
			newBoard.CandidateMap[len(line)] = append(existingCandidates, line)
		} else {
			newBoard.CandidateMap[len(line)] = []string{line}
		}
    }

	return newBoard, nil
}

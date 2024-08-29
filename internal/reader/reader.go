package reader

import (
	"bufio"
	"finnegan/internal/board"
	"os"
	"strconv"
)

func Read(gamePath string) (*board.Timeline, error) {
	// Open File
	file, err := os.Open(gamePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileReader := bufio.NewReader(file)
	newBoard := new(board.Board)
	timeline := new(board.Timeline)

	newBoard.DarkCell = '■'
	newBoard.WorkingRow = nil
	newBoard.WorkingCol = nil
	timeline.CandidateMap = make(map[int][]string)

	// Get Grid size
	gridSizeLine, err := fileReader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	boardSize, err := strconv.Atoi(gridSizeLine[:len(gridSizeLine)-1])
	if err != nil {
		return nil, err
	}

	// Initilize Grid
	newBoard.Grid = make([][]rune, boardSize)
	for i := range newBoard.Grid {
		newBoard.Grid[i] = make([]rune, boardSize)
	}

	// Gather grid from file
	for row := range boardSize {
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
		timeline.CandidateReference = append(timeline.CandidateReference, line)

		if existingCandidates, ok := timeline.CandidateMap[len(line)]; ok {
			timeline.CandidateMap[len(line)] = append(existingCandidates, line)
		} else {
			timeline.CandidateMap[len(line)] = []string{line}
		}
	}

	timeline.Boards = append(timeline.Boards, *newBoard)
	timeline.CurrentBoard = *newBoard
	timeline.Length = len(timeline.Boards)
	return timeline, nil
}

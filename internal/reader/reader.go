package reader

import (
	"finnegan/internal/board"
	"fmt"
	"os"
	"bufio"
)

func Read(gamePath string) (*board.Board, error) {
	file, err := os.Open(gamePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

package main

import (
	"finnegan/internal/reader"
	"fmt"
	"log"
)

func main() {
	board, err := reader.Read("./games/game1.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(board.Grid)
	for row := range board.Grid {
		for _, num := range board.Grid[row] {
			fmt.Print(string(num))
		}
		fmt.Println()
	}
	fmt.Println(board.NumberList)
}

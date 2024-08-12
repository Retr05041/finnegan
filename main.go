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

	fmt.Println(board.CandidateMap)
	if board.Solve() {
		board.Display()
	} else {
		fmt.Println("Could not solve puzzle.")
	}
}

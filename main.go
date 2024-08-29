package main

import (
	"finnegan/internal/reader"
	"finnegan/internal/solver"
	"fmt"
	"log"
)

func main() {
	board, err := reader.Read("./games/game1.txt")
	if err != nil {
		log.Fatal(err)
	}

	if solver.Solve(board) {
		fmt.Println("SOLVED")
		solver.CurrentBoard.Display()
	} else {
		fmt.Println("Could not solve puzzle.")
	}
}

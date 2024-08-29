package main

import (
	"finnegan/internal/reader"
	"finnegan/internal/solver"
	"fmt"
	"log"
)

func main() {
	timeline, err := reader.Read("./games/game1.txt")
	if err != nil {
		log.Fatal(err)
	}

	if solver.Solve(timeline) {
		fmt.Println("SOLVED")
		timeline.CurrentBoard.Display()
	} else {
		fmt.Println("Could not solve puzzle.")
	}
}

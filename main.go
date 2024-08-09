package main

import (
	"finnegan/internal/reader"
	"log"
)

func main() {
	_, err := reader.Read("./games/game1.txt")
	if err != nil {
		log.Fatal(err)
	}
}

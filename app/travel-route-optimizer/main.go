package main

import (
	"flag"
	"fmt"

	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	csvName = flag.String("routes", "./challenge_description/input-file.txt", "travel routes file")
)

func main() {
	flag.Parse() // `go run main.go -h` for help flag

	rots := database.Routes{}
	rots.LoadCsv(*csvName)

	for origin, options := range rots {
		for destination, price := range options {
			fmt.Println(origin, destination, price)
		}
	}
	// Up Server
}

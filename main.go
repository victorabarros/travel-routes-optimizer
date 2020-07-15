package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/victorabarros/challenge-bexs/app/server"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	csvName = flag.String("routes", "./challenge_description/input-file.txt", "travel routes file")
)

func main() {
	flag.Parse() // `go run main.go -h` for help flag

	rots := database.Routes{}
	if err := rots.LoadCsv(*csvName); err != nil {
		panic(err)
	}

	for origin, options := range rots {
		for destination, price := range options {
			fmt.Printf("%s\t%s\t%0.2f\n", origin, destination, price)
		}
	}

	// Up Server
	server.Run()

	// terminal
	time.Sleep(1 * time.Hour)
}

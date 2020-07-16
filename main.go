package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/victorabarros/challenge-bexs/app/server"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	csvName = flag.String("routes", "./input-file.txt", "travel routes file")
)

func main() {
	fmt.Println("Starting Service")
	flag.Parse() // `go run main.go -h` for help flag

	rots, err := database.New(*csvName)
	if err != nil {
		panic(err)
	}
	defer rots.File.Close()

	// Up Server
	server.Run(rots, *csvName)
	// terminal
	time.Sleep(1 * time.Hour) // apagar
}

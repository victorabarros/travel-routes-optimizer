package main

import (
	"flag"
	"fmt"

	"github.com/victorabarros/challenge-bexs/app/server"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

const (
	port = "8092" // TODO Move port to cfg
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
	server.Run(rots, *csvName, port)
}

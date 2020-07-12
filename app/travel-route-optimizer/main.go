package main

import (
	"flag"
	"fmt"
)

var (
	csvName = flag.String("routes", "./challenge_description/input-file.txt", "travel routes file")
)

func main() {
	flag.Parse() // `go run main.go -h` for help flag

	rots := database.routes{}
	rots.loadCsv(*csvName)

	for origin, options := range rots {
		for destiny, price := range options {
			fmt.Println(origin, destiny, price)
		}
	}
}

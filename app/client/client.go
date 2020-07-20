package main

// /home/vabarros/go/src/github.com/victorabarros/golang-class-notes/BlueBook_DonovanKernighan/4_Composite_Types/3/main.go

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/victorabarros/challenge-bexs/internal/database"
)

const (
	port = "8092"
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter the route")
	for input.Scan() {
		line := input.Text()
		if line != "" {
			handlerInput(line)
			fmt.Println("\nplease enter the route")
		}
	}
}

func handlerInput(line string) {
	lineSplited := strings.Split(line, "-")
	if len(lineSplited) != 2 {
		fmt.Println("invalid format, please enter \"ORG-DES\" format")
		return
	}

	orig, dest := lineSplited[0], lineSplited[1]
	url := fmt.Sprintf("http://challenge-bexs-server:%s/routes?origin=%s&destination=%s", port, orig, dest)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("server error")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(http.StatusText(resp.StatusCode))
		return
	}

	payload := struct {
		Amout    float64                `json:"amout"`
		Schedule map[int]database.Route `json:"schedule"`
	}{}
	json.NewDecoder(resp.Body).Decode(&payload)

	fmt.Print("best route: ")
	for _, rout := range payload.Schedule {
		fmt.Printf("%s - ", rout.Origin)
	}
	fmt.Printf("%s > $%.2f\n", dest, payload.Amout)
}

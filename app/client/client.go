package main

// /home/vabarros/go/src/github.com/victorabarros/golang-class-notes/BlueBook_DonovanKernighan/4_Composite_Types/3/main.go

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/victorabarros/challenge-bexs/internal/config"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	port string
	srv  string
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Error in load Enviromnts variables.")
	}
	port = strconv.Itoa(cfg.Port)
	srv = cfg.Server

	input := bufio.NewScanner(os.Stdin)
	msg := "please enter the route"
	fmt.Println(msg)

	for input.Scan() {
		line := input.Text()
		if line != "" {
			handlerInput(line)
			fmt.Println("\n" + msg)
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
	url := fmt.Sprintf("http://%s:%s/routes?origin=%s&destination=%s", srv, port, orig, dest)

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

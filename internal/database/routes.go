package database

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Routes it's simple dictionary of travel routes
// r[origin][destination] = price
type Routes map[string]map[string]int

// LoadCsv loads the startup csv file
func (r Routes) LoadCsv(csvName string) {
	fmt.Printf("Loading file %s\n", csvName)
	csvLines, err := readCsv(csvName)
	if err != nil {
		panic(err) // TODO better os.fatal or logrus.fatal
	}

	done := make(chan error)
	go func() {
		done <- r.fillRoutes(csvLines)
	}()

	for {
		for _, r := range `-\|/` {
			select {
			case <-done:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r%c", r)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (r Routes) fillRoutes(lines [][]string) (err error) {
	for idx, line := range lines {
		orig, dest := line[0], line[1]
		price, ok := strconv.Atoi(line[2])
		if ok != nil {
			fmt.Printf("Field \"%s\", from line %d, isn't a valid value for price\n", line[2], idx+1)
			panic("") // TODO better os.fatal or logrus.fatal
		}

		_, prs := r[orig]
		if !prs {
			r[orig] = make(map[string]int)
		}

		r[orig][dest] = price
	}
	return nil
}

// TODO: faz sentido a camada de DB abrir e um arquivo?
func readCsv(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

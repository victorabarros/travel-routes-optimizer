package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/victorabarros/challenge-bexs/internal"
)

// Routes it's simple dictionary of travel routes
// r[origin][destination] = price
type Routes map[string]map[string]float64

// LoadCsv loads the startup csv file
func (r Routes) LoadCsv(csvName string) error {
	fmt.Printf("Loading file %s\n", csvName)
	csvLines, err := internal.ReadCsv(csvName)
	if err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- r.fillRoutes(csvLines)
	}()

	for {
		select {
		case err := <-done:
			fmt.Print("\r")
			return err
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c", r)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}
}

func (r Routes) fillRoutes(lines [][]string) (err error) {
	for idx, line := range lines {
		orig, dest := line[0], line[1]
		price, ok := strconv.ParseFloat(line[2], 64)
		if ok != nil {
			return fmt.Errorf("field \"%s\", from line %d, isn't a valid value for price", line[2], idx+1)
		}

		_, prs := r[orig]
		if !prs {
			r[orig] = make(map[string]float64)
		}

		r[orig][dest] = price
	}
	return nil
}

package database

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// RouteDB is the model of travel routes at database layer
type RouteDB struct {
	db      map[string]map[string]float64 // TODO improte name
	csvName string
	File    *os.File
}

// New returns a new Route
func New(csvPath string) (RouteDB, error) {
	r := RouteDB{
		make(map[string]map[string]float64),
		csvPath,
		nil,
	}

	f, err := os.OpenFile(csvPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return r, err
	}
	r.File = f

	err = r.loadCsv()

	return r, err
}

// loadCsv loads the startup csv file
func (r RouteDB) loadCsv() error {
	fmt.Printf("Loading file %s\n", r.csvName)
	csvLines, err := csv.NewReader(r.File).ReadAll()
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

func (r RouteDB) fillRoutes(lines [][]string) (err error) {
	for idx, line := range lines {
		// Usar goroutine e mutex aqui
		orig, dest := line[0], line[1]
		price, ok := strconv.ParseFloat(line[2], 64)
		if ok != nil {
			return fmt.Errorf("field \"%s\", from line %d, isn't a valid value for price", line[2], idx+1)
		}

		_, prs := r.db[orig]
		if !prs {
			r.db[orig] = make(map[string]float64)
		}

		r.db[orig][dest] = price
	}
	return nil
}

// Route model
type Route struct {
	Origin      string
	Destination string
	Price       float64
}

// InsertRoute add new route to db
func (r RouteDB) InsertRoute(route Route) error {
	_, prs := r.db[route.Origin]
	if !prs {
		r.db[route.Origin] = make(map[string]float64)
	}
	r.db[route.Origin][route.Destination] = route.Price

	// TODO antes de escrever tem que valida se a linha já não existe usando o map.
	// Se já existir, como sobreescrever?
	if _, err := r.File.WriteString(
		fmt.Sprintf("%s,%s,%v\n", route.Origin, route.Destination, route.Price)); err != nil {
		return err
	}
	return nil
}

// PrintAll print all routes
func (r RouteDB) PrintAll() {
	// TODO remove
	fmt.Println("r.PrintAll")
	for origin, options := range r.db {
		for destination, price := range options {
			fmt.Printf("\n%s\t%s\t%0.2f\n", origin, destination, price)
		}
	}
}

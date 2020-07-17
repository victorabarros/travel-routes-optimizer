package database

import (
	"encoding/csv"
	"fmt"
	"math"
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

		_, prs := r.db[orig] // TODO usar upper
		if !prs {
			r.db[orig] = make(map[string]float64)
		}

		r.db[orig][dest] = price // TODO usar upper
	}
	return nil
}

// Route model
type Route struct { //Todo improve name. better transfer?
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Price       float64 `json:"price"`
}

// InsertRoute add new route to db
func (r RouteDB) InsertRoute(route Route) error {
	_, prs := r.db[route.Origin]
	if !prs {
		r.db[route.Origin] = make(map[string]float64) // TODO usar upper
	}
	r.db[route.Origin][route.Destination] = route.Price // TODO usar upper

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

// FindBestOffer find the cheapest transfer
func (r RouteDB) FindBestOffer(orig, dest string) (map[int]Route, float64) {
	var bestOffer float64 = math.MaxFloat64
	schedule := make(map[int]Route)

	_, prs := r.db[orig]
	if !prs || orig == dest {
		return schedule, bestOffer
	}

	return r.iterSearch(orig, dest, make(map[int]Route), make(map[int]Route), 0.0, bestOffer)
}

func (r RouteDB) iterSearch(orig, dest string,
	hist, bestSched map[int]Route,
	bugget, bestOffer float64) (map[int]Route, float64) {
	connxs, _ := r.db[orig]
	for connx, price := range connxs {
		_, prs := r.db[connx]
		if connx == dest && bugget+price < bestOffer {
			bestSched = make(map[int]Route)
			bestSched = deepCopy(hist)
			bestSched[len(bestSched)] = Route{
				orig,
				connx,
				price,
			}
			bestOffer = bugget + price
		} else if prs {
			_hist := make(map[int]Route)
			_hist = deepCopy(hist)
			_hist[len(_hist)] = Route{
				orig,
				connx,
				price,
			}
			bestSched, bestOffer = r.iterSearch(connx, dest, _hist,
				bestSched, bugget+price, bestOffer)
		}
	}
	return bestSched, bestOffer
}

func deepCopy(in map[int]Route) map[int]Route {
	out := make(map[int]Route)
	for k, v := range in {
		out[k] = v
	}
	return out
}

package database

import (
	"fmt"
	"os"
	"testing"
)

const (
	validPath   string = "./testSuccess.csv"
	invalidPath string = "./testFail.csv"
)

var (
	validRoutes RouteDB
	validLines  = []string{
		"GRU,BRC,10.5",
		"BRC,SCL,5",
		"GRU,CDG,75",
	}
	invalidLines = []string{
		"GRU,BRC,r4",
	}
)

func init() {
	// Write valid file
	if err := writeCsvFile(validPath, validLines); err != nil {
		panic(err)
	}

	valid, _ := New(validPath)
	validRoutes = valid

	if err := writeCsvFile(invalidPath, invalidLines); err != nil {
		panic(err)
	}
}

func writeCsvFile(path string, lines []string) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	for _, line := range lines {
		f.WriteString(line + "\n")
	}
	return nil
}

func TestLoadCsvSuccess(t *testing.T) {
	// defer os.Remove(validPath)
	// validRoutes.loadCsv(validPath)
	contains := func(sl []string, elem string) bool {
		for _, n := range sl {
			if elem == n {
				return true
			}
		}
		return false
	}

	for orig, v := range validRoutes.db {
		for dest, pric := range v {
			join := fmt.Sprintf("%s,%s,%v", orig, dest, pric)
			if !contains(validLines, join) {
				t.Errorf("\"%s\" doesn't contains \"%s\"", validLines, join)
			}
		}
	}
}

func TestLoadCsvInexistentFail(t *testing.T) {
	_, err := New("./inexistent.csv")
	if err == nil {
		t.Error("Must fail to read inexistent file")
	}
}

func TestLoadCsvInvalid(t *testing.T) {
	// defer os.Remove(invalidPath)
	_, err := New(invalidPath)
	if err == nil {
		t.Error("Must fail to load invalid price cell")
	}
}

func TestInsertRouteSuccess(t *testing.T) {
	if err := validRoutes.InsertRoute(Route{"A", "B", 12.4}); err != nil {
		t.Error(err)
	}
}

// func TestInsertRouteFail(t *testing.T) {
// 	path := "./empty.csv"
// 	if err := writeCsvFile(path, []string{"a,b,3"}); err != nil {
// 		t.Error(err)
// 	}
// 	r, err := New(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	fmt.Println("sleeping")
// 	time.Sleep(5 * time.Second)
// 	os.Remove(path)
// 	if err := r.InsertRoute(Route{"A", "B", 12.4}); err == nil {
// 		t.Error("Must fail to insert on inexistent file")
// 	}
// }

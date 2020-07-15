package database

import (
	"fmt"
	"os"
	"testing"
)

const (
	testFile string = "./testSuccess.csv"
)

var (
	testLines = []string{
		"GRU,BRC,10.5",
		"BRC,SCL,5",
		"GRU,CDG,75",
	}
)

func init() {
	if err := writeCsvFile(testFile, testLines); err != nil {
		panic(err)
	}
}

func writeCsvFile(name string, lines []string) error {
	f, err := os.Create(name)
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
	defer os.Remove(testFile)
	r := Routes{}
	r.LoadCsv(testFile)

	contains := func(sl []string, elem string) bool {
		for _, n := range sl {
			if elem == n {
				return true
			}
		}
		return false
	}

	for orig, v := range r {
		for dest, pric := range v {
			join := fmt.Sprintf("%s,%s,%v", orig, dest, pric)
			if !contains(testLines, join) {
				t.Errorf("\"%s\" doesn't contains \"%s\"", testLines, join)
			}
		}
	}
}

func TestLoadCsvInexistentFail(t *testing.T) {
	r := Routes{}
	if err := r.LoadCsv("./inexistent.csv"); err == nil {
		t.Error("Must fail to read inexistent file")
	}
}

func TestLoadCsvInvalid(t *testing.T) {
	name := "./fail.csv"
	writeCsvFile(name, []string{"a,b,3r"})
	defer os.Remove(name)
	r := Routes{}
	if err := r.LoadCsv(name); err == nil {
		t.Error("Must fail to load invalid price cell")
	}
}

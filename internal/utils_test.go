package internal

import (
	"os"
	"strings"
	"testing"
)

const (
	testFile string = "./testSuccess.csv"
)

var (
	testLines = []string{
		"a,b",
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

func TestReadCsvSucess(t *testing.T) {
	defer os.Remove(testFile)
	csv, err := ReadCsv(testFile)
	if err != nil {
		t.Error(err)
	}

	for idx, cells := range csv {
		if testLines[idx] != strings.Join(cells, ",") {
			t.Errorf("%s and %s must be equals", testLines[idx], strings.Join(cells, ","))
		}
	}
}

func TestReadCsvInvalidFile(t *testing.T) {
	name := "./inexistent.csv"
	_, err := ReadCsv(name)
	if err == nil {
		t.Error("Must fail to read inexistent file")
	}
}

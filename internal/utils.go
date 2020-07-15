package internal

import (
	"encoding/csv"
	"os"
)

// ReadCsv returns the cells from CSV
func ReadCsv(fileName string) ([][]string, error) {
	// TODO: Este cara em "internal" está no lugar certo?
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	lines, _ := csv.NewReader(file).ReadAll()

	return lines, nil
}

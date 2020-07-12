package database

import (
	"encoding/csv"
	"os"
	"strconv"
)

type routes map[string]map[string]int

func (r routes) loadCsv(csvName string) {
	csvFile, err := os.Open(csvName)
	defer csvFile.Close()
	if err != nil {
		panic(err)
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}

	for _, line := range csvLines {
		origin, destiny := line[0], line[1]
		price, ok := strconv.Atoi(line[2])
		if ok != nil {
			panic("strcon error")
			// TODO improve message
		}

		_, prs := r[origin]
		if !prs {
			r[origin] = make(map[string]int)
		}

		r[origin][destiny] = price
	}
}

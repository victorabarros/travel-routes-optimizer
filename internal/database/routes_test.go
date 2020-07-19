package database

import (
	"fmt"
	"math"
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
		"BRC,SCL,5",
		"GRU,BRC,10",
		"GRU,CDG,75",
		"GRU,ORL,56",
		"GRU,SCL,20",
		"ORL,CDG,5",
		"SCL,ORL,20",
	}
	invalidLines = []string{
		"GRU,BRC,r4",
	}
)

func init() {
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

func TestFindBestOfferSucess(t *testing.T) {
	cases := []struct {
		orig   string
		dest   string
		sched  map[int]Route
		amount float64
	}{
		{
			"BRC",
			"SCL",
			map[int]Route{
				0: Route{
					"BRC",
					"SCL",
					5,
				},
			},
			5,
		},
		{
			"GRU",
			"BRC",
			map[int]Route{
				0: Route{
					"GRU",
					"BRC",
					10,
				},
			},
			10,
		},
		{
			"ORL",
			"CDG",
			map[int]Route{
				0: Route{
					"ORL",
					"CDG",
					5,
				},
			},
			5,
		},
		{
			"SCL",
			"ORL",
			map[int]Route{
				0: Route{
					"SCL",
					"ORL",
					20,
				},
			},
			20,
		},
		{
			"BRC",
			"ORL",
			map[int]Route{
				0: Route{
					"BRC",
					"SCL",
					5,
				},
				1: Route{
					"SCL",
					"ORL",
					20,
				},
			},
			25,
		},
		{
			"GRU",
			"SCL",
			map[int]Route{
				0: Route{
					"GRU",
					"BRC",
					10,
				},
				1: Route{
					"BRC",
					"SCL",
					5,
				},
			},
			15,
		},
		{
			"SCL",
			"CDG",
			map[int]Route{
				0: Route{
					"SCL",
					"ORL",
					20,
				},
				1: Route{
					"ORL",
					"CDG",
					5,
				},
			},
			25,
		},
		{
			"BRC",
			"CDG",
			map[int]Route{
				0: Route{
					"BRC",
					"SCL",
					5,
				},
				1: Route{
					"SCL",
					"ORL",
					20,
				},
				2: Route{
					"ORL",
					"CDG",
					5,
				},
			},
			30,
		},
		{
			"GRU",
			"ORL",
			map[int]Route{
				0: Route{
					"GRU",
					"BRC",
					10,
				},
				1: Route{
					"BRC",
					"SCL",
					5,
				},
				2: Route{
					"SCL",
					"ORL",
					20,
				},
			},
			35,
		},
		{
			"GRU",
			"CDG",
			map[int]Route{
				0: Route{
					"GRU",
					"BRC",
					10,
				},
				1: Route{
					"BRC",
					"SCL",
					5,
				},
				2: Route{
					"SCL",
					"ORL",
					20,
				},
				3: Route{
					"ORL",
					"CDG",
					5,
				},
			},
			40,
		},
	}

	for _, _case := range cases {
		sched, amount := validRoutes.SearchBestOffer(_case.orig, _case.dest)
		if amount != _case.amount {
			t.Errorf("Amount %.2f and %.2f should be equals", amount, _case.amount)
		}
		if err := compareScheds(_case.sched, sched); err != nil {
			t.Error(err)
		}

		for idx, rout := range _case.sched {
			nextRout, prs := _case.sched[idx+1]
			if prs {
				if rout.Destination != nextRout.Origin {
					t.Errorf("route %d destination = %s should br equal origin to next route %d = %s | full sched: %+2v", idx, rout.Destination, idx+1, nextRout.Origin, _case.sched)
				}
			}
		}
	}
}

func compareScheds(sched1, sched2 map[int]Route) error {
	// TODO fazer um repo de funções utils com map[interface]interface
	// olhar no stack overflow
	var err error = nil

	if len(sched1) != len(sched2) {
		err = fmt.Errorf("%s\nschedules has differents length: %d and %d",
			err, len(sched1), len(sched2))
	}
	for k1, v1 := range sched1 {
		v2, prs := sched2[k1]
		if !prs {
			err = fmt.Errorf("%s\nkey %d doesn't exists at sched2",
				err, k1)
		}
		if v1 != v2 {
			err = fmt.Errorf("%s\nsched1[%d] = %+2v different from sched2[%d] = %+2v ",
				err, k1, v1, k1, v2)
		}
	}

	return err
}

func TestFindBestOfferNotFound(t *testing.T) {
	cases := []struct {
		orig   string
		dest   string
		sched  map[int]Route
		amount float64
	}{
		{
			"NotFound",
			"",
			nil,
			math.MaxFloat64,
		},
		{
			"GRU",
			"NotFound",
			nil,
			math.MaxFloat64,
		},
	}
	for _, _case := range cases {
		sched, amount := validRoutes.SearchBestOffer(_case.orig, _case.dest)
		if sched != nil || _case.amount != amount {
			t.Errorf("schedule %+2v must be nil", sched)
		}
	}
}

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/victorabarros/challenge-bexs/internal/database"
)

const (
	port             = "8094"
	validPath string = "./testSuccess.csv"
)

var (
	validRoutes database.RouteDB
	validLines  = []string{
		"BRC,SCL,5",
		"GRU,BRC,10",
		"GRU,CDG,75",
		"GRU,ORL,56",
		"GRU,SCL,20",
		"ORL,CDG,5",
		"SCL,ORL,20",
	}
)

func init() {
	setTests()
	go Run(validRoutes, validPath, port)
}

func setTests() {
	if err := writeCsvFile(validPath, validLines); err != nil {
		panic(err)
	}

	valid, _ := database.New(validPath)
	validRoutes = valid
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

func TestInsertSuccess(t *testing.T) {
	payload := database.Route{
		Origin:      "A",
		Destination: "B",
		Price:       0.3,
	}

	payloadDecoded, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Post("http://localhost:"+port+"/routes",
		"application/json", bytes.NewBuffer(payloadDecoded))
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("%+2v", json.NewDecoder(resp.Body))
	}
}

func TestInsertBadRequest1(t *testing.T) {
	payload := struct {
		field1 string
		field2 int
		field3 float64
		field4 bool
	}{
		"A",
		4,
		0.3,
		false,
	}

	payloadDecoded, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Post("http://localhost:"+port+"/routes",
		"application/json", bytes.NewBuffer(payloadDecoded))
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("%+2v", resp)
	}
}

func TestInsertBadRequest2(t *testing.T) {
	payloadRoute := []database.Route{
		{
			Origin:      "A",
			Destination: "",
			Price:       0.3,
		},
		{
			Origin:      "",
			Destination: "B",
			Price:       0.3,
		},
		{
			Origin:      "",
			Destination: "",
			Price:       0.3,
		},
	}

	for payload := range payloadRoute {
		payloadDecoded, err := json.Marshal(payload)
		if err != nil {
			t.Error(err)
		}

		resp, err := http.Post("http://localhost:"+port+"/routes",
			"application/json", bytes.NewBuffer(payloadDecoded))
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%+2v", json.NewDecoder(resp.Body))
		}
	}
}

func TestInsertError(t *testing.T) {
	payload := database.Route{
		Origin:      "A",
		Destination: "B",
		Price:       0.3,
	}

	payloadDecoded, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
	}

	os.Remove(validPath)
	resp, err := http.Post("http://localhost:"+port+"/routes",
		"application/json", bytes.NewBuffer(payloadDecoded))
	if err != nil {
		t.Error("Insert through inexistent file must fail")
	}

	if resp.StatusCode != http.StatusInternalServerError {
		fmt.Printf("%+2v", resp)
		t.Error("Insert through inexistent file must fail")
	}
	setTests()
}

func TestSearchSuccess(t *testing.T) {
	cases := []struct {
		Origin      string
		Destination string
		Amount      float64
	}{
		{
			"BRC",
			"SCL",
			5,
		},
		{
			"GRU",
			"BRC",
			10,
		},
		{
			"ORL",
			"CDG",
			5,
		},
		{
			"SCL",
			"ORL",
			20,
		},
		{
			"BRC",
			"ORL",
			25,
		},
		{
			"GRU",
			"SCL",
			15,
		},
		{
			"SCL",
			"CDG",
			25,
		},
		{
			"BRC",
			"CDG",
			30,
		},
		{
			"GRU",
			"ORL",
			35,
		},
		{
			"GRU",
			"CDG",
			40,
		},
	}

	baseURL := "http://localhost:" + port + "/routes"
	for _, _case := range cases {
		url := fmt.Sprintf("%s?origin=%s&destination=%s",
			baseURL, _case.Origin, _case.Destination)

		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}

		payload := struct {
			Amout    float64                `json:"amout"`
			Schedule map[int]database.Route `json:"schedule"`
		}{}

		json.NewDecoder(resp.Body).Decode(&payload)
		if payload.Amout != _case.Amount {
			t.Error("search cheapest schedule fail")
		}
	}
}

func TestSearchBadRequestFieldEmpty(t *testing.T) {
	cases := []struct {
		Origin      string
		Destination string
	}{
		{
			Origin:      "A",
			Destination: "",
		},
		{
			Origin:      "",
			Destination: "B",
		},
		{
			Origin:      "",
			Destination: "",
		},
	}

	baseURL := "http://localhost:" + port + "/routes"
	for _, _case := range cases {
		url := fmt.Sprintf("%s?origin=%s&destination=%s",
			baseURL, _case.Origin, _case.Destination)

		resp, err := http.Get(url)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			fmt.Printf("response: %+2v", resp)
			t.Errorf("must return bad request")
		}
	}

}

func TestSearchBadRequestNoField(t *testing.T) {
	baseURL := "http://localhost:" + port + "/routes"
	url := fmt.Sprintf("%s?origin=%s",
		baseURL, "Origin")

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}

	url = fmt.Sprintf("%s?destination=%s",
		baseURL, "Destination")

	resp, err = http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}
}

func TestSearchBadRequestFieldDuplicated(t *testing.T) {
	baseURL := "http://localhost:" + port + "/routes"
	url := fmt.Sprintf("%s?origin=%s&origin=%s&destination=%s&destination=%s",
		baseURL, "OriginA", "OriginB", "DestinationA", "DestinationB")

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}
}

func TestSearchNotFound(t *testing.T) {
	baseURL := "http://localhost:" + port + "/routes"
	url := fmt.Sprintf("%s?origin=%s&destination=%s",
		baseURL, "NOT", "B")

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNotFound {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}
}

func TestHealthCheckLiveness(t *testing.T) {
	url := "http://localhost:" + port + "/healthz"

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}
}

func TestHealthCheckReadness(t *testing.T) {
	url := "http://localhost:" + port + "/started"

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("response: %+2v", resp)
		t.Errorf("must return bad request")
	}
}

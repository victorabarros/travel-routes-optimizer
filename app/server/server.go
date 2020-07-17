package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

const (
	port = "8092" // TODO Move port to cfg
)

var (
	// cfg, _  = config.Load()
	db    database.RouteDB
	start = time.Now()
)

// Run up the server.
func Run(rout database.RouteDB, fileName string) {
	db = rout
	r := mux.NewRouter()
	r.HandleFunc("/routes", insert).Methods(http.MethodPost)
	r.HandleFunc("/routes", find).Methods(http.MethodGet)
	r.HandleFunc("/healthz", liveness)
	r.HandleFunc("/started", started)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	fmt.Printf("Up apllication at port %s\n", port)
	go panic(srv.ListenAndServe())
}

func insert(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Starting \"createBook\" route")
	payload := database.Route{}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(struct{ Message string }{err.Error()})
		return
	}

	if err := db.InsertRoute(payload); err != nil {
		panic(err)
	}

	rw.WriteHeader(http.StatusCreated)
}

func find(rw http.ResponseWriter, req *http.Request) {
	var err error = nil
	params := req.URL.Query()

	orig, prs := params["origin"]
	if !prs {
		err = fmt.Errorf("%s\n%s", err, "query param \"origin\" is required")
	} else if len(orig) > 1 {
		err = fmt.Errorf("%s\n%s", err, "only one query param \"origin\" is allowed")
	}

	dest, prs := params["destination"]
	if !prs {
		err = fmt.Errorf("%s\n%s", err, "query param \"destination\" is required")
	} else if len(dest) > 1 {
		err = fmt.Errorf("%s\n%s", err, "only one query param \"destination\" is allowed")
	}

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(struct{ Message string }{err.Error()})
		return
	}

	schedule, amount := db.FindBestOffer(orig[0], dest[0]) //TODO: Usar uppercase
	if len(schedule) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(struct {
		Amout    float64                `json:"amout"`
		Schedule map[int]database.Route `json:"schedule"`
	}{
		amount,
		schedule,
	})
}

// liveness is k8S liveness probe, returns if pod is alive
// Based on: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
func liveness(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(struct {
		ServiceName string
		Version     string
	}{
		"Travel-Route-Optimizer",
		"v0.1.0",
	})
}

// started returns how long is the container up
func started(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte((time.Since(start)).String()))
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	db    database.RouteDB
	start = time.Now()
)

// Run up the server.
func Run(rout database.RouteDB, fileName, port string) {
	db = rout
	r := mux.NewRouter()
	r.HandleFunc("/routes", insert).Methods(http.MethodPost)
	r.HandleFunc("/routes", search).Methods(http.MethodGet)
	r.HandleFunc("/healthz", liveness)
	r.HandleFunc("/started", started)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	logrus.Debug("Up apllication at port %s\n", port)
	panic(srv.ListenAndServe())
}

func insert(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"insert\" started")
	payload := database.Route{}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		logrus.Debug(err.Error())
		json.NewEncoder(rw).Encode(struct{ Message string }{err.Error()})
		return
	} else if payload.Origin == "" || payload.Destination == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(
			struct{ Message string }{"origin and destination can't be empty or nil"})
		return
	}

	if err := db.InsertRoute(payload); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "application/json")
		logrus.Error(err.Error())
		json.NewEncoder(rw).Encode(
			struct{ Message string }{http.StatusText(http.StatusInternalServerError)})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func search(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"search\" started")
	var err error = nil
	params := req.URL.Query()

	orig, prs := params["origin"]
	if !prs {
		err = fmt.Errorf("%s\n%s", err, "query param \"origin\" is required")
	} else if len(orig) != 1 {
		err = fmt.Errorf("%s\n%s", err, "only one query param \"origin\" is required")
	} else if orig[0] == "" {
		err = fmt.Errorf("%s\n%s", err, "\"origin\" is required")
	}

	dest, prs := params["destination"]
	if !prs {
		err = fmt.Errorf("%s\n%s", err, "query param \"destination\" is required")
	} else if len(dest) != 1 {
		err = fmt.Errorf("%s\n%s", err, "only one query param \"destination\" is required")
	} else if dest[0] == "" {
		err = fmt.Errorf("%s\n%s", err, "\"destination\" is required")
	}

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		logrus.Debug(err.Error())
		json.NewEncoder(rw).Encode(struct{ Message string }{err.Error()})
		return
	}

	sched, amount := db.SearchBestOffer(strings.ToUpper(orig[0]), strings.ToUpper(dest[0]))
	if len(sched) == 0 {
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
		sched,
	})
}

// liveness is k8S liveness probe, returns if pod is alive
// Based on: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
func liveness(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("route \"healthz\" started")
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
	logrus.Debug("route \"started\" started")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte((time.Since(start)).String()))
}

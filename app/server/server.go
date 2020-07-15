package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	port = "8092" // Move port to cfg
)

var (
	// cfg, _  = config.Load()
	started = time.Now()
)

// Run up the server.
func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", livenessController)
	r.HandleFunc("/started", startedController)
	r.HandleFunc("/routes", insertRoute).Methods(http.MethodPost)
	// r.HandleFunc("/routes", startedController).Methods(http.MethodGet) //  Consulta de melhor rota entre dois pontos.

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	fmt.Printf("Up apllication at port %s\n", port)
	go panic(srv.ListenAndServe())
}

func insertRoute(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Starting \"createBook\" route")
	payload := struct {
		Origin      string
		Destination string
		Price       float64
	}{}

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(struct{ Message string }{err.Error()})
		return
	}

	// Registra de novas rotas.Devem ser persistidas no arquivo csv
	// r[payload.Origin][payload.Destination] = payload.Price
	// f.WriteString(fmt.Sprintf("%s,%s,%v\n", payload.Origin, payload.Destination, payload.Price))

	rw.WriteHeader(http.StatusCreated)
}

// livenessController returns if pod is alive
// Based on: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
// k8S liveness probe
func livenessController(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Liveness route triggered")
	response := struct {
		ServiceName string
		Version     string
	}{
		"Travel-Transfer-Optimizer",
		"v0.1.0",
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

// startedController returns how long is the container up
func startedController(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Started route triggered")
	data := (time.Since(started)).String()

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(data))
}

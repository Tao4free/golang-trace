package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/gorilla/mux"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	gcp_project_id := os.Getenv("GCP_PROJECT_ID")
	exporter, err := stackdriver.NewExporter(
		stackdriver.Options{ProjectID: gcp_project_id})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample(),})

	router := registerHandler()

	var handler http.Handler = router
	handler = &ochttp.Handler{
		Handler:     handler,
		Propagation: &propagation.HTTPFormat{}}

	log.Fatalln(http.ListenAndServe(":8080", handler))
}

// Route handler
func registerHandler() *mux.Router {

	r := mux.NewRouter()

	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/hello", helloHandler).Methods("GET")
	s.HandleFunc("/again", againHandler).Methods("GET")

	return s
}

// Hello handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Acquire inbound context
	ctx := r.Context()
	time.Sleep(time.Second * 1)

	// Start span
	_, span := trace.StartSpan(ctx, "start=start_to_get_hello_string")
	defer span.End()
	h := generateHello()
	// Start span
	_, span_print := trace.StartSpan(ctx, "start=print_hello_string")
	defer span_print.End()

	fmt.Fprintf(w, h)
}

func generateHello() string {
	time.Sleep(time.Second * 1)
	h := "Hello?"
	return h
}

// Again handler
func againHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 1)
	fmt.Fprintf(w, "Again?")
}

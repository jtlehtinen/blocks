package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jtlehtinen/blocks/internal"
)

type ResponseEnvelope = map[string]any

func writeJSON(w http.ResponseWriter, status int, body any, header http.Header) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	for k, v := range header {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)

	return nil
}

func errorResponse(w http.ResponseWriter, status int, message any, header http.Header) {
	err := writeJSON(w, status, ResponseEnvelope{"error": message}, header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, allowedMethods []string) {
	header := http.Header{
		"Allow": []string{strings.Join(allowedMethods, ", ")},
	}

	message := fmt.Sprintf("method %s is not supported for this resource", r.Method)
	errorResponse(w, http.StatusMethodNotAllowed, message, header)
}

func blocksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowedResponse(w, r, []string{http.MethodGet})
		return
	}

	switch r.Method {
	case http.MethodGet:

	}
}

type application struct {
	chain *internal.Blockchain
}

func run() error {
	port := 3000
	router := http.NewServeMux()

	router.HandleFunc("/blocks", blocksHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return server.ListenAndServe()
}

func main() {
	/*
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	*/
	bc := internal.NewBlockchain()
	bc.AddBlock("foo")
	bc.AddBlock("bar")

	b, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(b))
}

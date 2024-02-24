package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func jsonResponse(response any, w http.ResponseWriter) {
	rJson, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		rJson = []byte(fmt.Sprintf("Failed to properly encode system data to JSON, with error: %s\n", err))
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rJson)
}

func middleware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileSystems, err := GetDisksMap()
		if err == nil {
			jsonResponse(fileSystems, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	})
}

func main() {
	host := flag.String("host", "localhost", "Port to run the server on")
	port := flag.Uint("port", 8080, "Port to run the server on")

	flag.Parse()

	fmt.Printf("Listening on %s:%d...\n", *host, *port)
	http.Handle("/", middleware())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}

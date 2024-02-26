package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const AppVersion = "0.1"

type response struct {
	Version string           `json:"version"`
	Disks   map[string]*disk `json:"disks"`
}

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

		response := response{
			Version: AppVersion,
			Disks:   fileSystems,
		}

		if err == nil {
			jsonResponse(response, w)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	})
}

func main() {
	host := flag.String("host", "localhost", "Port to run the server on")
	port := flag.Uint("port", 8930, "Port to run the server on")
	version := flag.Bool("version", false, "prints current app version")

	flag.Parse()
	if *version {
		fmt.Printf("swmonitor %s\n", AppVersion)
		os.Exit(0)
	}

	fmt.Printf("Listening on %s:%d...\n", *host, *port)
	http.Handle("/", middleware())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

const AppVersion = "0.2"

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
		disksMap, err := GetDisksMap()

		response := response{
			Version: AppVersion,
			Disks:   disksMap,
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

	fs := ff.NewFlagSet("gocanary")
	var (
		version = fs.BoolLong("version", "prints current app version")
		host    = fs.String('h', "host", "localhost", "Port to run the server on")
		port    = fs.Uint('p', "port", 8930, "Port to run the server on")
		_       = fs.String('c', "config", "", "Path to config file")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("GC"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("err=%v\n", err)
		os.Exit(0)
	}

	if *version {
		fmt.Printf("gocanary %s\n", AppVersion)
		os.Exit(0)
	}

	fmt.Printf("Listening on %s:%d...\n", *host, *port)
	http.Handle("/", middleware())
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}

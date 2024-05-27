package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Version string                  `json:"Version"`
	Alerts  *map[string]interface{} `json:"Alerts,omitempty"`
	Memory  map[string]*memoryStats `json:"Memory"`
	Disks   map[string]*disk        `json:"Disks"`
}

func writeJsonResponse(response response, w http.ResponseWriter) {
	rJson, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to properly encode system data to JSON, with error: %s\n", err)
		fmt.Println(errorMessage)
		rJson = []byte(errorMessage)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	responseCode := http.StatusOK
	if response.Alerts != nil {
		responseCode = http.StatusInternalServerError
	}
	w.WriteHeader(responseCode)
	w.Write(rJson)
}

func buildResponseHandler(alerts []alertRule) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		disksMap, err := getDisksMap()
		memoryMap, err := getMemoryMap()

		response := response{
			Version: AppVersion,
			Memory:  memoryMap,
			Disks:   disksMap,
		}

		if alerts := testAlertRules(response, alerts); len(alerts) > 0 {
			response.Alerts = &alerts
		}

		if err == nil {
			writeJsonResponse(response, w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	})
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Version string           `json:"version"`
	Disks   map[string]*disk `json:"disks"`
}

func writeJsonResponse(response any, w http.ResponseWriter) {
	rJson, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to properly encode system data to JSON, with error: %s\n", err)
		fmt.Println(errorMessage)
		rJson = []byte(errorMessage)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rJson)
}

func buildResponseHandler(alerts []alertRule) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		disksMap, err := getDisksMap()

		response := response{
			Version: AppVersion,
			Disks:   disksMap,
		}

		if alert, rule := testAlertRules(response, alerts); alert {
			fmt.Println(alert, *rule)
		}

		if err == nil {
			writeJsonResponse(response, w)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	})
}

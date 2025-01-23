package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWError(w http.ResponseWriter, code int, msg string, err error) {

	if err != nil {
		log.Println(err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWJson(w, code, errorResponse{Error: msg})
}

func respondWJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)

}

package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	const maxLen = 140
	forbiddenWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	type chirp struct {
		Body string `json:"body"`
	}
	type response struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err := decoder.Decode(&c)

	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	if len(c.Body) > maxLen {
		respondWError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	for _, word := range strings.Split(c.Body, " ") {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			c.Body = strings.Replace(c.Body, word, "****", -1)
		}
	}

	respondWJson(w, http.StatusOK, response{Cleaned_body: c.Body})

}

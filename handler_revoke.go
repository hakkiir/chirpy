package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {

	headers := req.Header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		respondWError(w, http.StatusBadRequest, "malformed authorization header", nil)
		return
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		respondWError(w, http.StatusBadRequest, "malformed authorization header", nil)
		return
	}

	refToken := splitAuth[1]

	err := cfg.db.RevokeToken(req.Context(), refToken)
	if err != nil {
		respondWError(w, http.StatusBadRequest, "refresh token not found", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

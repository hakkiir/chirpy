package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hakkiir/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerSetChirpyRed(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	type response struct {
		User
	}

	//check API key
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "API key not found", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWError(w, http.StatusUnauthorized, "API key does not match", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWError(w, http.StatusNoContent, "wrong event type", nil)
		return
	}

	parsedUUID, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "unable to parse user id", err)
		return
	}

	_, err = cfg.db.UpdateToChirpyRed(req.Context(), parsedUUID)
	if err != nil {
		respondWError(w, http.StatusNotFound, "user not found", err)
		return
	}

	respondWJson(w, http.StatusNoContent, response{})

}

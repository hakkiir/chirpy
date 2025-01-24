package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hakkiir/chirpy/internal/auth"
	"github.com/hakkiir/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "invalid bearer token", err)
		return
	}
	tokenUUID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "invalid bearer token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	//check crirp length
	const maxLen = 140
	if len(params.Body) > maxLen {
		respondWError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	//filter bad words
	forbiddenWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	for _, word := range strings.Split(params.Body, " ") {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			params.Body = strings.Replace(params.Body, word, "****", -1)
		}
	}

	//Insert chirp to database
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: tokenUUID,
	})
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error inserting data to database", err)
		return
	}
	//respond with JSON
	respondWJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

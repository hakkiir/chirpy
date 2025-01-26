package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error querying database", err)
		return
	}

	chirps := []Chirp{}

	for _, c := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}
	respondWJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, req *http.Request) {

	id := req.PathValue("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		respondWError(w, http.StatusBadRequest, "Error converting id to UUID", err)
		return
	}

	dbChirp, err := cfg.db.GetSingleChirp(req.Context(), uid)
	if err != nil {
		respondWError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	//respond with JSON
	respondWJson(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {

	nullUUID := uuid.NullUUID{
		UUID:  uuid.Nil,
		Valid: false,
	}

	//get author from query params, if not empty parse UUID and set to nullUUID struct
	authorId := req.URL.Query().Get("author_id")
	if authorId != "" {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			respondWError(w, http.StatusInternalServerError, "Error parsing uuid", err)
			return
		}
		nullUUID.UUID = authorUUID
		nullUUID.Valid = true
	}

	dbChirps, err := cfg.db.GetAllChirps(req.Context(), nullUUID)
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

	//sort if parameter given
	sortParam := req.URL.Query().Get("sort")
	if sortParam == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	}
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
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

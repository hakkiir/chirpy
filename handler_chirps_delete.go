package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hakkiir/chirpy/internal/auth"
	"github.com/hakkiir/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, req *http.Request) {

	id := req.PathValue("id")

	ChirpUid, err := uuid.Parse(id)
	if err != nil {
		respondWError(w, http.StatusBadRequest, "Error converting id to UUID", err)
		return
	}

	//get access token from header
	accessToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWError(w, http.StatusUnauthorized, "could not authorize", err)
		return
	}

	tokenUser, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "could not authorize", err)
		return
	}

	dbChirp, err := cfg.db.GetSingleChirp(req.Context(), ChirpUid)
	if err != nil {
		respondWError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	if dbChirp.UserID != tokenUser {
		respondWError(w, http.StatusForbidden, "you shall not delete", err)
		return
	}

	err = cfg.db.DeleteChirpById(req.Context(), database.DeleteChirpByIdParams{
		UserID: tokenUser,
		ID:     ChirpUid,
	})
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error deleting chirp:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

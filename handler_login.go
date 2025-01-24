package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hakkiir/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tokenExpirationTime := time.Hour

	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		tokenExpirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	JWT, err := auth.MakeJWT(user.ID, cfg.secret, tokenExpirationTime)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error creating token", err)
		return
	}

	respondWJson(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     JWT,
	})

}

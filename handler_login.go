package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hakkiir/chirpy/internal/auth"
	"github.com/hakkiir/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Decode request json
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	// Validate login
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

	// Create JWT access token

	const tokenExpirationTime = time.Hour
	JWT, err := auth.MakeJWT(user.ID, cfg.secret, tokenExpirationTime)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error creating access token", err)
		return
	}

	// Create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}
	_, err = cfg.db.InsertRefrestToken(req.Context(), database.InsertRefrestTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error inserting refresh token to db", err)
		return
	}

	respondWJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        JWT,
		RefreshToken: refreshToken,
	})

}

package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/hakkiir/chirpy/internal/auth"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(req.Context(), refToken)
	if err != nil {
		respondWError(w, http.StatusBadRequest, "reftoken not found", nil)
		return
	}
	if dbToken.RevokedAt.Valid || dbToken.ExpiresAt.Before(time.Now()) {
		respondWError(w, http.StatusUnauthorized, "refresh token expired or revoked", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "error creating access token", nil)
		return
	}
	respondWJson(w, http.StatusOK, response{
		Token: accessToken,
	})
}

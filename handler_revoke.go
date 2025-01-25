package main

import (
	"net/http"

	"github.com/hakkiir/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {

	refToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	err = cfg.db.RevokeToken(req.Context(), refToken)
	if err != nil {
		respondWError(w, http.StatusBadRequest, "refresh token not found", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

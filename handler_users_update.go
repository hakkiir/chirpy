package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hakkiir/chirpy/internal/auth"
	"github.com/hakkiir/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, req *http.Request) {

	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
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

	//decode request body
	decoder := json.NewDecoder(req.Body)
	r := reqBody{}
	err = decoder.Decode(&r)

	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	hashedPw, err := auth.HashPassword(r.Password)
	if err != nil {
		log.Print("error hashing pw")
		respondWError(w, http.StatusInternalServerError, "Error hashing pw", err)
		return
	}
	user, err := cfg.db.UpdateEmailAndPassword(req.Context(), database.UpdateEmailAndPasswordParams{
		ID:             tokenUser,
		Email:          r.Email,
		HashedPassword: hashedPw,
	})

	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error updating database", err)
		return
	}

	respondWJson(w, http.StatusOK, response{
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}})

}

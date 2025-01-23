package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {

	type reqBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	r := reqBody{}
	err := decoder.Decode(&r)

	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error decoding params", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), r.Email)
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error querying database", err)
		return
	}

	respondWJson(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

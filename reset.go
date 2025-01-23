package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		respondWError(w, http.StatusInternalServerError, "Error querying database", err)
		return
	}
	cfg.fileserverRequests.Store(0)
	w.WriteHeader(http.StatusOK)
}

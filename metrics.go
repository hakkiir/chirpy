package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverRequests.Add(0))
	b := []byte(s)
	req.Header.Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverRequests.Add(1)
		next.ServeHTTP(w, r)
	})
}

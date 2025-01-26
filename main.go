package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/hakkiir/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverRequests atomic.Int32
	db                 *database.Queries
	platform           string
	secret             string
}

func main() {

	//load .env
	godotenv.Load(".env")

	//get variables from .env
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	//open db connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("db connection failed")
	}

	apiCfg := apiConfig{
		fileserverRequests: atomic.Int32{},
		db:                 database.New(db),
		platform:           platform,
		secret:             secret,
	}

	const port = "8080"
	const filepathRoot = "app/"

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetSingleChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerChirpDelete)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

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
}

func main() {

	//load .env
	godotenv.Load(".env")

	//get dbURL from .env
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	//open db connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("db connection failed")
	}

	apiCfg := apiConfig{
		fileserverRequests: atomic.Int32{},
		db:                 database.New(db),
		platform:           platform,
	}

	const port = "8080"
	const filepathRoot = "app/"

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}

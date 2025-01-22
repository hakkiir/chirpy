package main

import (
	"log"
	"net/http"
)

const port = "8080"

func main() {

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

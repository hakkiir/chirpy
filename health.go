package main

import "net/http"

func handlerHealth(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

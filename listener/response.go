package main

import (
	"net/http"
)

func OKResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	// nolint
	w.Write([]byte("OK"))
}

func ErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

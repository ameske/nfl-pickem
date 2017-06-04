package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Logo(logosDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team := mux.Vars(r)["team"]
		file := fmt.Sprintf("%s/%s.gif", logosDir, team)
		http.ServeFile(w, r, file)
	}
}

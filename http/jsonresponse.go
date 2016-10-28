package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, response interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	err := enc.Encode(response)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}

type statusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func WriteJSONSuccess(w http.ResponseWriter, response string) {
	w.WriteHeader(200)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	err := enc.Encode(&statusResponse{Status: "success", Message: response})
	if err != nil {
		log.Println(err)
	}
}

func WriteJSONError(w http.ResponseWriter, code int, response string) {
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	err := enc.Encode(&statusResponse{Status: "failed", Message: response})
	if err != nil {
		log.Println(err)
	}
}

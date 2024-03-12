package services

import (
	"encoding/json"
	"log"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token != "" {
		resp := make(map[string]string, 1)
		resp["error"] = "Already authorized"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Marshal error: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}

		w.WriteHeader(http.StatusForbidden)
		w.Write(jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

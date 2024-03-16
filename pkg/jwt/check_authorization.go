package jwt

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckAuthorization(token string, path string) (int, []byte) {
	resp := make(map[string]string, 1)

	if token != "" && (path == "/sign-up" || path == "/sign-in" || path == "/login" || path == "/registration") {
		resp["error"] = "Already authorized"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Marshal error: " + err.Error())
			return http.StatusInternalServerError, []byte("error while marshal response")
		}

		return http.StatusForbidden, jsonResp
	} else if token == "" && !(path == "/sign-up" || path == "/sign-in" || path == "/login" || path == "/registration") {
		resp["error"] = "Not authorized"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Marshal error: " + err.Error())
			return http.StatusInternalServerError, []byte("error while marshal response")
		}

		return http.StatusForbidden, jsonResp
	}

	return http.StatusOK, nil
}

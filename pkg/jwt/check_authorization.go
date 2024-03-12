package jwt

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckAuthorization(token string, path string) (int, []byte) {
	if token != "" && (path == "/sign-up" || path == "/sign-in" || path == "/login" || path == "/registration") {
		resp := make(map[string]string, 1)
		resp["error"] = "Already authorized"

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Marshal error: " + err.Error())
			return http.StatusInternalServerError, []byte("Internal Server Error")
		}

		return http.StatusForbidden, jsonResp
	}

	return http.StatusOK, nil
}

package jwt

import (
	"filmoteka/internal/models"
	"net/http"
)

func CheckAuthorization(token string, path string) (int, []models.ErrorInfo) {
	resp := make(map[string]string, 1)
	var errorsPull []models.ErrorInfo

	if token != "" && (path == "/sign-up" || path == "/sign-in" || path == "/login" || path == "/registration") {
		resp["error"] = "Already authorized"

		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "Already authorized"})
		// jsonResp, err := json.Marshal(resp)
		// if err != nil {
		// 	log.Println("Marshal error: " + err.Error())
		// 	errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while marshal response"})
		// 	return http.StatusInternalServerError, errorsPull
		// }

		return http.StatusForbidden, errorsPull
	} else if token == "" && !(path == "/sign-up" || path == "/sign-in" || path == "/login" || path == "/registration") {
		resp["error"] = "Not authorized"

		errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "Not authorized"})
		// jsonResp, err := json.Marshal(resp)
		// if err != nil {
		// 	log.Println("Marshal error: " + err.Error())
		// 	errorsPull = append(errorsPull, models.ErrorInfo{Type: "error", Message: "error while marshal response"})
		// 	return http.StatusInternalServerError, errorsPull
		// }

		return http.StatusUnauthorized, errorsPull
	}

	return http.StatusOK, nil
}

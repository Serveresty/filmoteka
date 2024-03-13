package services

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Authorization")
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

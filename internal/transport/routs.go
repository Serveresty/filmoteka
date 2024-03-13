package transport

import (
	"filmoteka/internal/services"
	"net/http"
)

func Routes(mux *http.ServeMux) {
	//GET
	mux.HandleFunc("/sign-up", services.SignUp)
	mux.HandleFunc("/sign-in", services.SignIn)

	//POST
	mux.HandleFunc("/registration", services.Registration)
	mux.HandleFunc("/login", services.Login)
	mux.HandleFunc("/logout", services.Logout)
}

package transport

import (
	"filmoteka/internal/services"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Routes(mux *http.ServeMux) {
	//GET Auth
	mux.HandleFunc("/sign-up", services.SignUp)
	mux.HandleFunc("/sign-in", services.SignIn)

	//POST Auth
	mux.HandleFunc("/registration", services.Registration)
	mux.HandleFunc("/login", services.Login)
	mux.HandleFunc("/logout", services.Logout)

	/*----------------------------------------------------------------*/

	//GET film
	mux.HandleFunc("/films", services.GetFilms)
	mux.HandleFunc("/films/search", services.GetFilmsByFilter)

	//POST film
	mux.HandleFunc("/new-film", services.AddNewFilm)

	//PUT film
	mux.HandleFunc("/edit-film", services.EditInfoFilm)

	//DELETE film
	mux.HandleFunc("/delete-film", services.DeleteFilm)

	/*----------------------------------------------------------------*/

	//GET actor
	mux.HandleFunc("/actors", services.GetActors)

	//POST actor
	mux.HandleFunc("/new-actor", services.AddNewActor)

	//PUT actor
	mux.HandleFunc("/edit-actor", services.EditInfoActor)

	//DELETE film
	mux.HandleFunc("/delete-actor", services.DeleteActor)

	//SWAGGER
	mux.HandleFunc("/spec", handleSwaggerFile())
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/spec"),
	))
}

func handleSwaggerFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../docs/swagger.json")
	}
}

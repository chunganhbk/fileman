package http

import (
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/raedahgroup/fileman/config"
	"github.com/raedahgroup/fileman/storage"
	"github.com/rs/cors"
	"net/http"
)

type modifyRequest struct {
	What  string   `json:"what"`  // Answer to: what data type?
	Which []string `json:"which"` // Answer to: which fields?
}

func NewHandler(storage *storage.Storage, config config.ConfigState) (http.Handler, error) {
	r := mux.NewRouter()

	monkey := func(fn handleFunc) http.Handler {
		return handle(fn, storage, config)
	}

	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/login", monkey(loginHandler)).Methods("POST")
	api.Handle("/signup", monkey(signupHandler))
	api.Handle("/renew", monkey(renewHandler))

	users := api.PathPrefix("/users").Subrouter()
	users.Handle("", monkey(usersGetHandler)).Methods("GET")
	users.Handle("", monkey(userPostHandler)).Methods("POST")
	users.Handle("/{id:[0-9]+}", monkey(userPutHandler)).Methods("PUT")
	users.Handle("/{id:[0-9]+}", monkey(userGetHandler)).Methods("GET")
	users.Handle("/{id:[0-9]+}", monkey(userDeleteHandler)).Methods("DELETE")

	api.PathPrefix("/raw").Handler(monkey(rawHandler)).Methods("GET")

	resources := api.PathPrefix("/resources").Subrouter()
	resources.PathPrefix("/").Handler(monkey(resourceGetHandler)).Methods("GET")
	resources.PathPrefix("/").Handler(monkey(resourceDeleteHandler)).Methods("DELETE")
	resources.PathPrefix("/").Handler(monkey(resourcePostPutHandler)).Methods("POST")
	resources.PathPrefix("/").Handler(monkey(resourcePostPutHandler)).Methods("PUT")
	resources.PathPrefix("/").Handler(monkey(resourcePatchHandler)).Methods("PATCH")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(rice.MustFindBox("../web/dist").HTTPBox())))
	//return http.StripPrefix(config.BaseURL, r), nil
	return c.Handler(r), nil
}


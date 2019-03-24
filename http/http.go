package http

import (
	"encoding/json"
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/raedahgroup/fileman/config"
	"github.com/raedahgroup/fileman/storage"
	"net/http"
	"strings"
	"text/template"
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
	index, static := getStaticHandlers()
	r.NotFoundHandler = index
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", static))
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
	resources.Use(FindForder)
	resources.PathPrefix("/").Handler(monkey(resourceGetHandler)).Methods("GET")
	resources.PathPrefix("/").Handler(monkey(resourceDeleteHandler)).Methods("DELETE")
	resources.PathPrefix("/").Handler(monkey(resourcePostPutHandler)).Methods("POST")
	resources.PathPrefix("/").Handler(monkey(resourcePostPutHandler)).Methods("PUT")
	resources.PathPrefix("/").Handler(monkey(resourcePatchHandler)).Methods("PATCH")

	/*c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})*/

	return http.StripPrefix(config.BaseURL, r), nil
	//return c.Handler(r), nil
}
func getStaticHandlers() (http.Handler, http.Handler) {
	box := rice.MustFindBox("../web/dist")
	handler := http.FileServer(box.HTTPBox())
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r);
		}
		w.Header().Set("x-xss-protection", "1; mode=block")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		staticURL := strings.TrimPrefix(config.State.BaseURL +"/static", "/")
		data := map[string]interface{}{
			"BASE_URL":         "",
			"NAME":  "File manger",
			"StaticURL":       staticURL,
			"LOCALE":  "en",
		}
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data["Json"] = string(b)
		emberView :=  template.Must(template.New("index").Delims("[{[", "]}]").Parse(box.MustString("index.html")))
		if err := emberView.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	});
	static := handler
	return index, static
}
// replace path  api/resources
func FindForder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		r.URL.Path = strings.Replace(r.URL.Path, "api/resources/", "", -1)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

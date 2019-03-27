package http

import (
	"fmt"
	"github.com/raedahgroup/fileman/config"
	"github.com/raedahgroup/fileman/storage"
	"github.com/raedahgroup/fileman/users"
	"log"
	"net/http"
	"strconv"
)

type handleFunc func(w http.ResponseWriter, r *http.Request, d *data) (int, error)

type data struct {

	store    *storage.Storage
	config   *config.Server
	user     *users.User
	raw      interface{}
}

// Check implements rules.Checker.
func (d *data) Check(path string) bool {
	for _, rule := range d.user.Rules {
		if rule.Matches(path) {
			return rule.Allow
		}
	}

	return true
}

func handle(fn handleFunc, storage *storage.Storage, config *config.Server) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		status, err := fn(w, r, &data{
			store:  storage,
			config: config,
		})

		if status != 0 {
			txt := http.StatusText(status)
			http.Error(w, strconv.Itoa(status)+" "+txt, status)
		}

		if status >= 400 || err != nil {
			log.Printf("%s: %v %s %v", r.URL.Path, status, r.RemoteAddr, err)
		}
		fmt.Println("trang thai status %s", status, err)
	})

	//return http.StripPrefix(prefix, handler)
	return  handler
}

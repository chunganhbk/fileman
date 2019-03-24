package http

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/raedahgroup/fileman/errors"
	"github.com/raedahgroup/fileman/users"
	"net/http"
	"strings"
	"time"
)

type userInfo struct {
	ID           uint              `json:"id"`
	Locale       string            `json:"locale"`
	ViewMode     users.ViewMode    `json:"viewMode"`
	Name         string    `json:"name"`
	Perm         users.Permissions `json:"perm"`
	LockPassword bool              `json:"lockPassword"`
}

type authToken struct {
	User userInfo `json:"user"`
	jwt.StandardClaims
}

type extractor []string

func (e extractor) ExtractToken(r *http.Request) (string, error) {
	token, _ := request.HeaderExtractor{"Authorization"}.ExtractToken(r)
	// Checks if the token isn't empty and if it contains two dots.
	// The former prevents incompatibility with URLs that previously
	// used basic auth.
	if token != "" && strings.Count(token, ".") == 2 {
		splitToken := strings.Split(token, "Bearer ")
		reqToken := splitToken[1]
		return reqToken, nil
	}

	auth := r.URL.Query().Get("auth")
	if auth == "" {
		return "", request.ErrNoTokenInRequest
	}

	return auth, nil
}

func withUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return []byte(d.config.JWTKEY), nil
		}

		var tk authToken
		token, err := request.ParseFromRequestWithClaims(r, &extractor{}, &tk, keyFunc)
		if err != nil || !token.Valid {
			return http.StatusForbidden, nil
		}

		expired := !tk.VerifyExpiresAt(time.Now().Add(time.Hour).Unix(), true)
		updated := d.store.Users.LastUpdate(tk.User.ID) > tk.IssuedAt

		if expired || updated {
			w.Header().Add("X-Renew-Token", "true")
		}

		d.user, err = d.store.Users.Get(d.config.RootPath, tk.User.ID)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		return fn(w, r, d)
	}
}

func withAdmin(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if !d.user.Perm.Admin {
			return http.StatusForbidden, nil
		}

		return fn(w, r, d)
	})
}
var loginHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	if r.Body == nil {
		return http.StatusForbidden, nil
	}
	var cred signupBody;
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		return http.StatusForbidden, nil
	}
	user, err :=  d.store.Users.Get(d.config.RootPath, cred.Username)
	fmt.Println(user, err)
	if err != nil || !users.CheckPwd(cred.Password, user.Password) {
		return http.StatusForbidden, nil
	} else {
		return printToken(w, r, d, user)
	}
}

type signupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var signupHandler = func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	if r.Body == nil {
		return http.StatusBadRequest, nil
	}

	info := &signupBody{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if info.Password == "" || info.Username == "" {
		return http.StatusBadRequest, nil
	}

	user := &users.User{
		Username: info.Username,
	}

	pwd, err := users.HashPwd(info.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user.Password = pwd
	err = d.store.Users.Save(user)
	if err == errors.ErrExist {
		return http.StatusConflict, err
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

var renewHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	return printToken(w, r, d, d.user)
})

func printToken(w http.ResponseWriter, r *http.Request, d *data, user *users.User) (int, error) {
	claims := &authToken{
		User: userInfo{
			ID:           user.ID,
			Name:         string(user.Username[0]),
			Locale:       user.Locale,
			ViewMode:     user.ViewMode,
			Perm:         user.Perm,
			LockPassword: user.LockPassword,
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
			Issuer:    "File Manger",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(d.config.JWTKEY))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "cty")
	w.Write([]byte(signed))
	return 0, nil
}

package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/go-zoo/bone"
	"github.com/google/go-github/github"
	"github.com/naoina/toml"
	"github.com/pocke/gh-feeds/db"
)

type Config struct {
	Port         int
	CleintID     string
	ClientSecret string
}

var config *Config

var oauthConf = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/autorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func setConfig() error {
	f, err := Asset("config/application.toml")
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(f, config); err != nil {
		return err
	}

	oauthConf.ClientID = config.CleintID
	oauthConf.ClientSecret = config.ClientSecret
	return nil
}

func main() {
	if err := setConfig(); err != nil {
		panic(err)
	}

	mux := bone.New()
	mux.GetFunc("/authorize", Auth)
	mux.GetFunc("/auth_callback", Callback)

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	// TODO: state
	u := oauthConf.AuthCodeURL("state")
	http.Redirect(w, r, u, http.StatusMovedPermanently)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")
	tok, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghc := github.NewClient(oauthConf.Client(oauth2.NoContext, tok))
	u, _, err := ghc.Users.Get("") // get me
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := &db.User{
		ID:   *u.ID,
		Name: *u.Login,
		Auth: tok.AccessToken,
	}

	if _, err := user.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

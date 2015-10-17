package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/go-zoo/bone"
	"github.com/google/go-github/github"
	"github.com/naoina/toml"
	"github.com/pocke/gh-feeds/db"
	"github.com/pocke/hlog"
)

type Config struct {
	Port         int
	ClientID     string `toml:"client_id"`
	ClientSecret string
}

var config = new(Config)

var oauthConf = &oauth2.Config{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
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

	oauthConf.ClientID = config.ClientID
	oauthConf.ClientSecret = config.ClientSecret
	return nil
}

func main() {
	if err := db.UseProd(); err != nil {
		panic(err)
	}
	db.LogMode(true)

	if err := setConfig(); err != nil {
		panic(err)
	}

	mux := bone.New()
	mux.GetFunc("/authorize", Auth)
	mux.GetFunc("/auth_callback", Callback)

	log.Printf("Start HTTP Server with port %d", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), hlog.Wrap(mux.ServeHTTP))
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
	_, err = db.CreateUser(&db.UserParams{
		ID:   *u.ID,
		Name: *u.Login,
		Auth: tok.AccessToken,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

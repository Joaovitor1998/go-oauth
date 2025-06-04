package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	config "github.com/Joaovitor1998/language-social-auth"
	"github.com/Joaovitor1998/language-social-auth/utils"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load("../.env.development")
	if err != nil {
		utils.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	gConf := config.NewOAuthGoogleConfig(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_ENDPOINT"),
	)

	gAuthCodeConfig := config.OAuthCodeURLConfig{ State: "state-token", Opts: []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}}

	google := config.NewOAuthGoogle(gConf, gAuthCodeConfig)
	gf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(r.Context().Value(config.OAuthGoogleUserInfo{}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	mux.HandleFunc("/auth/google/login", google.InitiateLogin)
	mux.Handle("/auth/google/callback", google.Callback(gf))

	fConf := config.NewOAuthGoogleConfig(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_ENDPOINT"),
	)

	fAuthCodeConfig := config.OAuthCodeURLConfig{ State: "state-token", Opts: []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}}

	facebook := config.NewOAuthGoogle(fConf, fAuthCodeConfig)
	ff := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(r.Context().Value(config.OAuthFacebookUserInfo{}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	mux.HandleFunc("/auth/facebook/login", facebook.InitiateLogin)
	mux.Handle("/auth/facebook/callback", facebook.Callback(ff))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	
	host := fmt.Sprintf(":%s", os.Getenv("INTERNAL_PORT"))
	
	http.ListenAndServe(host, mux)
}
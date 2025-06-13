package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	logger "github.com/Joaovitor1998/go-logger"
	config "github.com/Joaovitor1998/go-oauth"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load("../.env.development")
	if err != nil {
		logger.Error("could not load .env file")
	}

	mux := http.NewServeMux()
	
	// GOOGLE OAUTH
	gConf := config.NewOAuthGoogleConfig(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_ENDPOINT"),
	)
	google, gCb := googleOauth(gConf)
	

	mux.HandleFunc("/auth/google/login", google.InitiateLogin)
	mux.Handle("/auth/google/callback", google.Callback(gCb))


	// FACEBOOK OAUTH
	fConf := config.NewOAuthGoogleConfig(
		os.Getenv("FACEBOOK_CLIENT_ID"),
		os.Getenv("FACEBOOK_CLIENT_SECRET"),
		os.Getenv("FACEBOOK_REDIRECT_ENDPOINT"),
	)
	facebook, fCb := facebookOauth(fConf)

	mux.HandleFunc("/auth/facebook/login", facebook.InitiateLogin)
	mux.Handle("/auth/facebook/callback", facebook.Callback(fCb))
	
	host := fmt.Sprintf(":%s", os.Getenv("INTERNAL_PORT"))
	http.ListenAndServe(host, mux)
}

func googleOauth(gConf *oauth2.Config) (google config.OAuthConfig, handler http.HandlerFunc) {
	// GOOGLE OAUTH
	gAuthCodeConfig := config.OAuthCodeURLConfig{ State: "state-token", Opts: []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}}

	google = config.NewOAuthGoogle(gConf, gAuthCodeConfig)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(r.Context().Value(config.OAuthGoogleUserInfo{}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	return google, handler
}

func facebookOauth(fConf *oauth2.Config) (facebook config.OAuthConfig, handler http.HandlerFunc){
	fAuthCodeConfig := config.OAuthCodeURLConfig{ State: "state-token", Opts: []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}}

	facebook = config.NewOAuthGoogle(fConf, fAuthCodeConfig)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(r.Context().Value(config.OAuthFacebookUserInfo{}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	return facebook, handler
}
package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type oAuthFacebook struct {
	config *oauth2.Config
	authCodeConfig OAuthCodeURLConfig
}

var fb_user_info_url = "https://graph.facebook.com/me?fields=id,name,email,first_name,last_name,picture&access_token"

func NewOAuthFacebookConfig(id, secret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: id,
		ClientSecret: secret,
		RedirectURL: redirectURL,
		Endpoint: facebook.Endpoint,
		Scopes: []string{"email", "public_profile"},
	}
}

func NewOAuthFacebook(config *oauth2.Config, authCodeConfig OAuthCodeURLConfig) *oAuthFacebook {
	return &oAuthFacebook{
		config: config,
		authCodeConfig: authCodeConfig,
	}
}

func (o *oAuthFacebook) InitiateLogin(w http.ResponseWriter, r *http.Request) {
	url := o.config.AuthCodeURL(o.authCodeConfig.State, o.authCodeConfig.Opts...)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (o *oAuthFacebook) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Verify state parameter
		state := r.URL.Query().Get("state")
		if state != o.authCodeConfig.State {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		// 2. Exchange code for token
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Authorization code not provided", http.StatusBadRequest)
			return
		}

		token, err := o.config.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}	

		// 3. Get user info from Facebook
		// Facebook's Graph API endpoint		
		client := o.config.Client(r.Context(), token)
		resp, err := client.Get(fb_user_info_url)
		if err != nil {
			http.Error(w, "failed getting user info: "+err.Error(), http.StatusBadRequest)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, "failed reading response body:" + string(body), http.StatusBadRequest)
			return
		}

		var userInfo OAuthFacebookUserInfo
		err = json.NewDecoder(resp.Body).Decode(&userInfo)
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), OAuthFacebookUserInfo{}, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type OAuthFacebookUserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}
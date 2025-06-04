package config

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oAuthGoogle struct {
	config *oauth2.Config
	authCodeConfig OAuthCodeURLConfig
}

var gg_user_info_url = "https://www.googleapis.com/oauth2/v2/userinfo?access_token"

// It creates the auth code URL and then it redirects you to login and calls the callback
func (o *oAuthGoogle) InitiateLogin(w http.ResponseWriter, r *http.Request) {
	url := o.config.AuthCodeURL(o.authCodeConfig.State, o.authCodeConfig.Opts...)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// It's a middleware that gets the user info and add its
// data into the request as `OAuthGoogleUserInfo`
func (o *oAuthGoogle) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			http.Error(w, "Authorization code not provided", http.StatusBadRequest)
			return
		}
		
		token, err := o.config.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token: " + err.Error(), http.StatusBadRequest)
			return
		}
		
		client := o.config.Client(r.Context(), token)
		resp, err := client.Get(gg_user_info_url)
		if err != nil {
			http.Error(w, "Failed to get user info: " + err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, "Google API error: " + string(body), resp.StatusCode)
			return
		}
		
		var userInfo OAuthGoogleUserInfo

		err = json.NewDecoder(resp.Body).Decode(&userInfo)
		if err != nil {
			http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), OAuthGoogleUserInfo{}, userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

    // Then create session as in LoginWithCredentials
	// 7. Find or create user in your database (implementation depends on your storage)
	// user, err := findOrCreateUserFromGoogle(userInfo)
	// if err != nil {
	// 	http.Error(w, "Failed to process user: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// 8. Create session or JWT token (implementation depends on your auth system)
	// err = createUserSession(w, r, user)
	// if err != nil {
	// 	http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// 9. Redirect to your application's frontend or return success
	// http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// 6. Parse the user info response
type OAuthGoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	VerifiedEmail bool `json:"verified_email"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture string `json:"picture"`
}

func NewOAuthGoogleConfig(id, secret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:   redirectURL,
		Endpoint: google.Endpoint,
		Scopes: []string {
		"https://www.googleapis.com/auth/userinfo.profile", 
		"https://www.googleapis.com/auth/userinfo.email",
		},
	}
}

func NewOAuthGoogle(config *oauth2.Config, authCodeConfig OAuthCodeURLConfig) *oAuthGoogle {
	return &oAuthGoogle{
		config: config,
		authCodeConfig: authCodeConfig,
	}
}
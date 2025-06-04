package config

import (
	"net/http"

	"golang.org/x/oauth2"
)

type OAuthConfig interface {
	InitiateLogin(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}

type OAuthCodeURLConfig struct {
	State string
	Opts []oauth2.AuthCodeOption
}

// Helper function to find or create user
// func findOrCreateUserFromGoogle(googleUser struct {
// 	ID      string `json:"id"`
// 	Email   string `json:"email"`
// 	Name    string `json:"name"`
// 	Picture string `json:"picture"`
// }) (*models.User, error) {
// 	// Implementation depends on your database
// 	// Typically you would:
// 	// 1. Check if user exists with this Google ID
// 	// 2. If not, check if user exists with this email
// 	// 3. If not, create new user
// 	// 4. Return the user object
// 	return &models.User{}, nil
// }

// // Helper function to create session
// func createUserSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
// 	// Implementation depends on your session management
// 	// Could be:
// 	// 1. Cookie-based session
// 	// 2. JWT token in cookie or response body
// 	return nil
// }
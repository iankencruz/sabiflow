package auth

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"), // e.g. http://localhost:8080/api/auth/google/callback
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func GetGoogleAuthURL(state string) string {
	return googleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return googleOAuthConfig.Exchange(ctx, code)
}

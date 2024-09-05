package utility

import (
	"mohhefni/go-blog-app/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ConfigGoogle(config config.OAuthConfig) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  config.ClientCallbackUrl,
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

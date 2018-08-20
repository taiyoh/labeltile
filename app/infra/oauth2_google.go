package infra

import (
	"github.com/taiyoh/labeltile/app"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuth2Google struct {
	app.OAuth2Google
	oauth2 *oauth2.Config
}

func NewOAuth2Google(client_id, client_secret, redirect_url string) *OAuth2Google {
	c := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  redirect_url,
		Scopes: []string{
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}
	return &OAuth2Google{oauth2: c}
}

func (o *OAuth2Google) AuthCodeURL(state string) string {
	return o.oauth2.AuthCodeURL(state)
}

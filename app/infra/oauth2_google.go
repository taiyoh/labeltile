package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/taiyoh/labeltile/app"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth2Google provides google oauth2 interface for getting user's profile
type OAuth2Google struct {
	app.OAuth2Google
	oauth2Conf  *oauth2.Config
	userInfoURL string
}

type tokenInfo struct {
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Gender        string `json:"gender"`
	Email         string `json:"email"`
	Link          string `json:"link"`
	GivenName     string `json:"given_name"`
	ID            string `json:"id"`
	VerifiedEmail bool   `json:"verified_email"`
}

// OAuth2GoogleTokenInfo is data transfer object for user profile
type OAuth2GoogleTokenInfo struct {
	app.OAuth2GoogleTokenInfo
	email     string
	expiresAt time.Time
	userID    string
}

// Email returns user's email from profile
func (i *OAuth2GoogleTokenInfo) Email() string {
	return i.email
}

// ExpiresAt returns when user token expires
func (i *OAuth2GoogleTokenInfo) ExpiresAt() time.Time {
	return i.expiresAt
}

// UserID returns user's id from profile
func (i *OAuth2GoogleTokenInfo) UserID() string {
	return i.userID
}

// NewOAuth2Google returns new OAuth2Google object
func NewOAuth2Google(clientID, clientSecret, redirectURL string) *OAuth2Google {
	c := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"profile",
			"email",
		},
		Endpoint: google.Endpoint,
	}
	return &OAuth2Google{
		oauth2Conf:  c,
		userInfoURL: "https://www.googleapis.com/userinfo/v2/me",
	}
}

// AuthCodeURL returns URL to oauth permission page
func (o *OAuth2Google) AuthCodeURL(state string) string {
	return o.oauth2Conf.AuthCodeURL(state)
}

func NewUserTokenInfo(id, email string, exp time.Time) *OAuth2GoogleTokenInfo {
	return &OAuth2GoogleTokenInfo{
		userID:    id,
		email:     email,
		expiresAt: exp,
	}
}

// GetTokenInfo returns OAuth2GoogleTokenInfo object using oauth2 protocol
func (o *OAuth2Google) GetTokenInfo(code string) (app.OAuth2GoogleTokenInfo, error) {
	ctx := context.Background()
	t, err := o.oauth2Conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	if !t.Valid() {
		return nil, errors.New("invalid token")
	}

	client := o.oauth2Conf.Client(ctx, t)

	resp, err := client.Get(o.userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var resBody tokenInfo
	b := bytes.NewBuffer([]byte{})
	b.ReadFrom(resp.Body)
	if err := json.Unmarshal(b.Bytes(), &resBody); err != nil {
		return nil, errors.New("broken json returns")
	}

	return NewUserTokenInfo(resBody.ID, resBody.Email, t.Expiry), nil
}

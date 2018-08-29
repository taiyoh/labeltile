package mock

import (
	"github.com/taiyoh/labeltile/app"
)

type OAuth2Google struct {
	app.OAuth2Google
	err  error
	info app.OAuth2GoogleTokenInfo
}

func (o *OAuth2Google) FillResponse(info app.OAuth2GoogleTokenInfo, err error) {
	o.info = info
	o.err = err
}

func (o *OAuth2Google) AuthCodeURL(state string) string {
	return ""
}

func (o *OAuth2Google) GetTokenInfo(code string) (app.OAuth2GoogleTokenInfo, error) {
	return o.info, o.err
}

func LoadOAuth2Google() *OAuth2Google {
	return &OAuth2Google{}
}

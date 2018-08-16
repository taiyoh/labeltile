package labeltile

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/taiyoh/labeltile/app/infra"
)

type requestUser struct {
	ID         string
	expireDate string
}

// Request provides variable, query and user grouping per request
type Request struct {
	Variables map[string]interface{}
	Query     string
	User      *requestUser
}

// NewRequest returns Request object with json and token validation
func NewRequest(body io.ReadCloser, userToken string, s *infra.UserTokenSerializer) (*Request, error) {
	r := &Request{}
	if err := json.NewDecoder(body).Decode(r); err != nil {
		return nil, errors.New("broken request")
	}
	if r.Variables == nil || r.Query == "" {
		return nil, errors.New("require query and variables")
	}
	if userToken == "" {
		return r, nil
	}

	if claims, err := s.Deserialize(userToken); err == nil {
		r.User = &requestUser{
			ID:         claims["userID"].(string),
			expireDate: claims["expireDate"].(string),
		}
		return r, nil
	}

	return nil, errors.New("broken user token")
}

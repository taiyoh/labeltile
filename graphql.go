package labeltile

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/resolver"
)

type requestUser struct {
	ID         string
	expireDate string
}

// GraphQL provides variable, query and user grouping per request for GraphQL execution
type GraphQL struct {
	Variables map[string]interface{}
	Query     string
	User      *requestUser
}

// NewGraphQLRequest returns Request object with json and token validation
func NewGraphQLRequest(body io.ReadCloser, userToken string, serializer app.UserTokenSerializer) (*GraphQL, error) {
	r := &GraphQL{}
	if err := json.NewDecoder(body).Decode(r); err != nil {
		return nil, errors.New("broken request")
	}
	if r.Variables == nil || r.Query == "" {
		return nil, errors.New("require query and variables")
	}
	if userToken == "" {
		return r, nil
	}

	if claims, err := serializer.Deserialize(userToken); err == nil {
		r.User = &requestUser{
			ID:         claims["userID"].(string),
			expireDate: claims["expireDate"].(string),
		}
		return r, nil
	}

	return nil, errors.New("broken user token")
}

func (g *GraphQL) Run(ctx context.Context) interface{} {
	ctx = context.WithValue(ctx, resolver.RequestUserCtxKey, g.User)
	r := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  g.Query,
		VariableValues: g.Variables,
		Context:        ctx,
	})
	return r.Data
}

var (
	schema graphql.Schema
)

func init() {
	resolver.InitializeTypes()
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: resolver.GetType(resolver.GQLType("RootQuery")),
	})
}

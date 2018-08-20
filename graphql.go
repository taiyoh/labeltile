package labeltile

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/resolver"
)

// GraphQL provides variable, query and user grouping per request for GraphQL execution
type GraphQL struct {
	Variables map[string]interface{}
	Query     string
	User      app.RequestUser
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
		userID := claims["userID"].(string)
		expDate := claims["expireDate"].(string)
		ed, _ := strconv.ParseInt(expDate, 10, 64)
		r.User = app.NewRequestUser(userID, ed)
		return r, nil
	}

	return nil, errors.New("broken user token")
}

func (g *GraphQL) Run(container app.Container, ctx context.Context) interface{} {
	ctx = context.WithValue(ctx, app.RequestUserCtxKey, g.User)
	ctx = context.WithValue(ctx, app.ContainerCtxKey, container)
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

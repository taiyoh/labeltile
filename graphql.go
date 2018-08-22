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

// GraphQL provides variable, query and user grouping per request for GraphQL execution
type GraphQL struct {
	Variables map[string]interface{}
	Query     string
	Schema    *graphql.Schema
}

// NewGraphQLRequest returns Request object with json and token validation
func NewGraphQLRequest(schema *graphql.Schema, body io.ReadCloser) (*GraphQL, error) {
	r := &GraphQL{Schema: schema}
	if err := json.NewDecoder(body).Decode(r); err != nil {
		return nil, errors.New("broken request")
	}
	if r.Variables == nil || r.Query == "" {
		return nil, errors.New("require query and variables")
	}
	return r, nil
}

// Run provides GraphQL request with resolver
func (g *GraphQL) Run(ctx context.Context) map[string]interface{} {
	r := graphql.Do(graphql.Params{
		Schema:         *g.Schema,
		RequestString:  g.Query,
		VariableValues: g.Variables,
		Context:        ctx,
	})
	res := map[string]interface{}{
		"data": r.Data,
	}
	if r.HasErrors() {
		res["errors"] = r.Errors
	}
	return res
}

// InitializeGraphQLSchema provides type initialization for GraphQL
func InitializeGraphQLSchema(container app.Container) *graphql.Schema {
	s := resolver.InitializeTypes(container)
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: s.Get(resolver.GQLType("RootQuery")),
	})
	return &schema
}

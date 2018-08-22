package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

// GQLType is registered GraphQL object type
type GQLType string

var types map[GQLType]*graphql.Object

// GetType returns GraphQL object which is registered
func GetType(t GQLType) *graphql.Object {
	o, _ := types[t]
	return o
}

// InitializeTypes provides user definition GraphQL types initialization
func InitializeTypes(container app.Container) {
	types = map[GQLType]*graphql.Object{
		GQLType("User"): buildType("User", []*graphql.Field{
			&graphql.Field{Name: "id", Type: graphql.NewNonNull(graphql.ID)},
			&graphql.Field{Name: "mail", Type: graphql.String},
			&graphql.Field{Name: "roles", Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String)))},
		}),
		GQLType("RootQuery"):    buildType("RootQuery", []*graphql.Field{}),
		GQLType("RootMutation"): buildType("RootMutation", []*graphql.Field{}),
	}
	initRootQuery(container)
	initRootMutation(container)
}

func buildType(name string, fieldList []*graphql.Field) *graphql.Object {
	o := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: graphql.Fields{},
	})
	for _, f := range fieldList {
		o.AddFieldConfig(f.Name, f)
	}
	return o
}

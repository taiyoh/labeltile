package resolver

import (
	"github.com/graphql-go/graphql"
)

// GQLType is registered GraphQL object type
type GQLType string

// CtxKey is access key for context.Context
type CtxKey string

var types map[GQLType]*graphql.Object

// GetType returns GraphQL object which is registered
func GetType(t GQLType) *graphql.Object {
	o, _ := types[t]
	return o
}

// InitializeTypes provides user definition GraphQL types initialization
func InitializeTypes() {
	types = map[GQLType]*graphql.Object{
		GQLType("User"): buildType("User", []*graphql.Field{
			buildField("id", graphql.NewNonNull(graphql.ID)),
			buildField("mail", graphql.String),
			buildField("roles", graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String)))),
		}),
		GQLType("RootQuery"):    buildType("RootQuery", []*graphql.Field{}),
		GQLType("RootMutation"): buildType("RootMutation", []*graphql.Field{}),
	}
	initRootQuery()
	initRootMutation()
}

func buildField(name string, o graphql.Output) *graphql.Field {
	return &graphql.Field{
		Name: name,
		Type: o,
	}
}

func buildResolverField(name string, o graphql.Output, args graphql.FieldConfigArgument, r graphql.FieldResolveFn) *graphql.Field {
	return &graphql.Field{
		Name:    name,
		Type:    o,
		Args:    args,
		Resolve: r,
	}
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

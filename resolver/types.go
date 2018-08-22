package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

// GQLType is registered GraphQL object type
type GQLType string

// TypeStorage provides stores GraphQL type definitions
type TypeStorage struct {
	stores map[GQLType]*graphql.Object
}

// Get returns graphql.Object which matches by given name
func (s *TypeStorage) Get(n GQLType) *graphql.Object {
	o, _ := s.stores[n]
	return o
}

// Register provides building graphql.Object and setting it
func (s *TypeStorage) Register(name string, fieldList ...*graphql.Field) {
	o := graphql.NewObject(graphql.ObjectConfig{
		Name:   name,
		Fields: graphql.Fields{},
	})
	for _, f := range fieldList {
		o.AddFieldConfig(f.Name, f)
	}
	s.stores[GQLType(name)] = o
}

// InitializeTypes provides user definition GraphQL types initialization
func InitializeTypes(container app.Container) *TypeStorage {
	s := &TypeStorage{stores: map[GQLType]*graphql.Object{}}
	s.Register("User",
		&graphql.Field{Name: "id", Type: graphql.NewNonNull(graphql.ID)},
		&graphql.Field{Name: "mail", Type: graphql.String},
		&graphql.Field{Name: "roles", Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String)))},
	)
	s.Register("RootQuery")
	s.Register("RootMutation")
	s.InitRootQuery(container)
	s.InitRootMutation(container)

	return s
}

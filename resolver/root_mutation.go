package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

// InitRootMutation fills resolver fields in RootMutation type
func (s *TypeStorage) InitRootMutation(container app.Container) {
	m := s.Get(GQLType("RootMutation"))
	rm := &RootMutation{container: container}
	for _, f := range []*graphql.Field{
		&graphql.Field{
			Name: "updateUser",
			Type: s.Get(GQLType("User")),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: rm.UpdateUser,
		},
	} {
		m.AddFieldConfig(f.Name, f)
	}
}

// RootMutation is field resolver aggregation for RootMutation type
type RootMutation struct {
	container app.Container
}

// UpdateUser is implementation for "updateUser" field in RootMutation type
func (t *RootMutation) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

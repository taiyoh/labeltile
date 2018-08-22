package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

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

type RootMutation struct {
	container app.Container
}

func (t *RootMutation) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

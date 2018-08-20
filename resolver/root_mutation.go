package resolver

import (
	"github.com/graphql-go/graphql"
)

func initRootMutation() {
	m := GetType(GQLType("RootMutation"))
	for _, f := range []*graphql.Field{
		&graphql.Field{
			Name: "updateUser",
			Type: GetType(GQLType("User")),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: rootMutationUpdateUser,
		},
	} {
		m.AddFieldConfig(f.Name, f)
	}
}

func rootMutationUpdateUser(p graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

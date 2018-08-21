package resolver

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

func initRootQuery() {
	r := GetType(GQLType("RootQuery"))
	for _, f := range []*graphql.Field{
		&graphql.Field{
			Name:    "operator",
			Type:    GetType(GQLType("User")),
			Resolve: rootQueryOperator,
		},
	} {
		r.AddFieldConfig(f.Name, f)
	}
}

func rootQueryOperator(p graphql.ResolveParams) (interface{}, error) {
	container, cok := p.Context.Value(app.ContainerCtxKey).(app.Container)
	if !cok {
		return nil, errors.New("container not found")
	}
	userID, rok := p.Context.Value(app.UserIDCtxKey).(string)
	if !rok {
		return nil, nil
	}
	u := app.UserFindService(userID, container)
	return u, nil
}

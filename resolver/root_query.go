package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

func initRootQuery(container app.Container) {
	r := GetType(GQLType("RootQuery"))
	rq := &RootQuery{container: container}
	for _, f := range []*graphql.Field{
		&graphql.Field{
			Name:    "operator",
			Type:    GetType(GQLType("User")),
			Resolve: rq.Operator,
		},
	} {
		r.AddFieldConfig(f.Name, f)
	}
}

type RootQuery struct {
	container app.Container
}

func (t *RootQuery) Operator(p graphql.ResolveParams) (interface{}, error) {
	container := t.container
	userID, rok := p.Context.Value(app.UserIDCtxKey).(string)
	if !rok {
		return nil, nil
	}
	u := app.UserFindService(userID, container)
	return u, nil
}

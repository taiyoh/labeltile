package resolver

import (
	"github.com/graphql-go/graphql"
	"github.com/taiyoh/labeltile/app"
)

func initRootQuery(s *TypeStorage, container app.Container) {
	r := s.Get(GQLType("RootQuery"))
	rq := &RootQuery{container: container}
	for _, f := range []*graphql.Field{
		&graphql.Field{
			Name:    "operator",
			Type:    s.Get(GQLType("User")),
			Resolve: rq.Operator,
		},
	} {
		r.AddFieldConfig(f.Name, f)
	}
}

// RootQuery is field resolver aggregation for RootQuery type
type RootQuery struct {
	container app.Container
}

// Operator is implementation for "operator" field in RootQuery type
func (t *RootQuery) Operator(p graphql.ResolveParams) (interface{}, error) {
	container := t.container
	userID, rok := p.Context.Value(app.UserIDCtxKey).(string)
	if !rok {
		return nil, nil
	}
	u := app.UserFindService(userID, container)
	return u, nil
}

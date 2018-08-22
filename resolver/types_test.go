package resolver_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app/infra/mock"
	"github.com/taiyoh/labeltile/resolver"
)

func TestInitialize(t *testing.T) {
	c := mock.LoadContainer()
	s := resolver.InitializeTypes(c)
	if o := s.Get(resolver.GQLType("RootQuery")); o == nil {
		t.Error("RootQuery not initialized")
	}
	if o := s.Get(resolver.GQLType("RootMutation")); o == nil {
		t.Error("RootMutation not initialized")
	}
	if o := s.Get(resolver.GQLType("User")); o == nil {
		t.Error("User not initialized")
	}
}

package resolver_test

import (
	"testing"

	"github.com/taiyoh/labeltile/resolver"
)

func TestInitialize(t *testing.T) {
	resolver.InitializeTypes()
	if o := resolver.GetType(resolver.GQLType("RootQuery")); o == nil {
		t.Error("RootQuery not initialized")
	}
	if o := resolver.GetType(resolver.GQLType("RootMutation")); o == nil {
		t.Error("RootMutation not initialized")
	}
	if o := resolver.GetType(resolver.GQLType("User")); o == nil {
		t.Error("User not initialized")
	}
}

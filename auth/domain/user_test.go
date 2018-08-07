package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
)

func TestUser(t *testing.T) {
	u := domain.NewUser(domain.UserID("1"), domain.UserMail("foo@example.com"), []domain.UserRoleID{domain.UserRoleID("viewer")})
	if u.ID != domain.UserID("1") {
		t.Error("user id should be 1")
	}
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.UserRoleID("viewer") {
		t.Error("user role should be viewer")
	}

	u = u.AddRole(domain.UserRoleID("editor"))
	if len(u.Roles) != 2 {
		t.Error("user roles count should be 2")
	}
	if u.Roles[1] != domain.UserRoleID("editor") {
		t.Error("user role should be editor")
	}

	u = u.DeleteRole(domain.UserRoleID("viewer"))
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.UserRoleID("editor") {
		t.Error("user role should be only editor")
	}
}

package auth

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
)

func TestUserService(t *testing.T) {
	reg := NewMockRegistry()
	s := NewUserService(reg)

	s.Register("hoge@example.com", "reader")

	repo := reg.UserRepository()
	u := repo.Find(domain.UserID("user:1"))
	if u == nil {
		t.Error("user not registered")
	}
	if u.Mail != "hoge@example.com" {
		t.Error("wrong mail address registered")
	}
	if len(u.Roles) != 1 {
		t.Error("role count should be 1")
	}
	if u.Roles[0] != domain.UserRoleID("reader") {
		t.Error("wrong role is registered")
	}

	s.AddRole("user:1", "editor")
	u = repo.Find(domain.UserID("user:1"))
	if len(u.Roles) != 2 {
		t.Error("role count should be 2")
	}
	if u.Roles[1] != domain.UserRoleID("editor") {
		t.Error("wrong role is registered")
	}

	s.DeleteRole("user:1", "reader")
	u = repo.Find(domain.UserID("user:1"))
	if len(u.Roles) != 1 {
		t.Error("role count should be 1")
	}
	if u.Roles[0] != domain.UserRoleID("editor") {
		t.Error("wrong role is registered")
	}

}

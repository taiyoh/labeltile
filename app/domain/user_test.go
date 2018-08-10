package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestUser(t *testing.T) {
	repo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID("1")
	})
	factory := domain.NewUserFactory(repo)
	u := factory.Build(domain.UserMail("foo@example.com"))
	if u.ID != domain.UserID("1") {
		t.Error("user id should be 1")
	}
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.RoleViewer {
		t.Error("user role should be viewer")
	}

	u = u.AddRole(domain.RoleEditor)
	if len(u.Roles) != 2 {
		t.Error("user roles count should be 2")
	}
	u = u.AddRole(domain.RoleViewer)
	if len(u.Roles) != 2 {
		t.Error("user roles count should be 2")
	}

	if u.Roles[1] != domain.RoleEditor {
		t.Error("user role should be editor")
	}

	u = u.DeleteRole(domain.RoleViewer)
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.RoleEditor {
		t.Error("user role should be only editor")
	}
}

func TestSpecifyUserRegistration(t *testing.T) {
	repo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID("1")
	})

	s := domain.NewUserSpecification(repo)

	addr := "foo@example.com"

	if err := s.SpecifyUserRegistration("foo bar baz"); err == nil {
		t.Error("invalid address should returns error")
	}

	if err := s.SpecifyUserRegistration("foo bar baz <foo@example.com>"); err == nil {
		t.Error("contains somenting other than E-mail address")
	}

	factory := domain.NewUserFactory(repo)
	u := factory.Build(domain.UserMail(addr))
	repo.Save(u)

	if err := s.SpecifyUserRegistration(addr); err == nil {
		t.Error("already registered")
	}

	if err := s.SpecifyUserRegistration("bar@example.com"); err != nil {
		t.Error("unregistered address should not returns error")
	}
}

package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestUser(t *testing.T) {
	repo := &mock.UserRepositoryImpl{
		DispenseIDFunc: func() domain.UserID {
			return domain.UserID("1")
		},
	}
	factory := domain.NewUserFactory(repo)
	u := factory.Build(domain.UserMail("foo@example.com"))
	if u.ID != domain.UserID("1") {
		t.Error("user id should be 1")
	}
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.UserRoleViewer {
		t.Error("user role should be viewer")
	}

	u = u.AddRole(domain.UserRoleEditor)
	if len(u.Roles) != 2 {
		t.Error("user roles count should be 2")
	}
	if u.Roles[1] != domain.UserRoleEditor {
		t.Error("user role should be editor")
	}

	u = u.DeleteRole(domain.UserRoleViewer)
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != domain.UserRoleEditor {
		t.Error("user role should be only editor")
	}
}

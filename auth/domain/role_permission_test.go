package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestRoleSpecification(t *testing.T) {
	rrepo := &domain.RoleRepository{}
	userID := 1
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		suid := string(userID)
		userID++
		return domain.UserID(suid)
	})

	s := domain.NewRoleSpecification(rrepo)

	factory := domain.NewUserFactory(urepo)

	op := factory.Build(domain.UserMail("foo@example.com"))
	if err := s.SpecifyEditRole(op, []domain.RoleID{domain.RoleEditor}); err == nil {
		t.Error("not permitted")
	}

	op = op.AddRole(domain.RoleManageUser)
	if err := s.SpecifyEditRole(op, []domain.RoleID{domain.RoleViewer}); err == nil {
		t.Error("Viewer role can't edit")
	}
	if err := s.SpecifyEditRole(op, []domain.RoleID{domain.RoleEditor}); err != nil {
		t.Error("operator should have permission and roles should be valid")
	}
}

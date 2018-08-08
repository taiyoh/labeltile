package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestRoleSpecification(t *testing.T) {
	rrepo := &domain.RoleRepository{}
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID("1")
	})

	s := domain.NewRoleSpecification(rrepo)

	factory := domain.NewUserFactory(urepo)

	op := factory.Build(domain.UserMail("foo@example.com"))
	if err := s.SpecifyEditRole(op); err == nil {
		t.Error("not permitted")
	}

	op = op.AddRole(domain.RoleManageUser)
	if err := s.SpecifyEditRole(op); err != nil {
		t.Error("operator should have permission")
	}
}

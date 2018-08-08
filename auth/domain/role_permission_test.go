package domain_test

import (
	"strconv"
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

	roleEditorID := strconv.Itoa(int(domain.RoleEditor))

	op := factory.Build(domain.UserMail("foo@example.com"))
	if err := s.SpecifyEditRole(op, []string{roleEditorID}); err == nil {
		t.Error("not permitted")
	}

	op = op.AddRole(domain.RoleManageUser)
	if err := s.SpecifyEditRole(op, []string{"1111111111"}); err == nil {
		t.Error("given roles are invalid")
	}
	if err := s.SpecifyEditRole(op, []string{"!!!"}); err == nil {
		t.Error("given roles are invalid")
	}
	if err := s.SpecifyEditRole(op, []string{roleEditorID}); err != nil {
		t.Error("operator should have permission and roles should be valid")
	}
}

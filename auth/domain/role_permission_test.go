package domain_test

import (
	"strconv"
	"testing"

	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestConvertRoleToID(t *testing.T) {
	rrepo := &domain.RoleRepository{}
	s := domain.NewRoleSpecification(rrepo)

	if _, err := s.ConvertRoleToID([]string{}); err == nil {
		t.Error("require role list")
	}

	roleEditor := strconv.Itoa(int(domain.RoleEditor))
	roleManager := strconv.Itoa(int(domain.RoleManageUser))

	if _, err := s.ConvertRoleToID([]string{"!!"}); err == nil {
		t.Error("invalid role given")
	}

	if _, err := s.ConvertRoleToID([]string{roleEditor, "1234567890"}); err == nil {
		t.Error("invalid role given")
	}

	if roleIDs, err := s.ConvertRoleToID([]string{roleEditor, roleManager}); err != nil {
		t.Error("all roles converted")
	} else {
		if len(roleIDs) != 2 {
			t.Error("roles are given for 2")
		}
		if roleIDs[0] != domain.RoleEditor {
			t.Error("given first is RoleEditor")
		}
		if roleIDs[1] != domain.RoleManageUser {
			t.Error("given second is RoleManager")
		}
	}
}

func TestSpecificationForAddingRole(t *testing.T) {
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
	if err := s.SpecifyAddRole(op, op, []domain.RoleID{domain.RoleEditor}); err == nil {
		t.Error("not permitted")
	}

	op = op.AddRole(domain.RoleManageUser)
	if err := s.SpecifyAddRole(op, op, []domain.RoleID{domain.RoleViewer}); err == nil {
		t.Error("Viewer role can't edit")
	}
	if err := s.SpecifyAddRole(op, op, []domain.RoleID{domain.RoleEditor}); err != nil {
		t.Error("operator should have permission and roles should be valid")
	}
}

package auth_test

import (
	"strconv"
	"testing"

	"github.com/taiyoh/labeltile/auth"
	"github.com/taiyoh/labeltile/auth/domain"
	"github.com/taiyoh/labeltile/auth/infra/mock"
)

func TestUserRegisterService(t *testing.T) {
	userID := 1
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		suid := string(userID)
		userID++
		return domain.UserID(suid)
	})
	f := domain.NewUserFactory(urepo)
	op := f.Build(domain.UserMail("foo@example.com"))
	urepo.Save(op)

	rrepo := &domain.RoleRepository{}

	if err := auth.UserRegisterService(string(op.ID), "target@example.com", urepo, rrepo); err == nil {
		t.Error("operator has no permission")
	}

	op = op.AddRole(domain.RoleManageUser)
	urepo.Save(op)

	if err := auth.UserRegisterService(string(op.ID), "foo@example.com", urepo, rrepo); err == nil {
		t.Error("already registered")
	}

	if err := auth.UserRegisterService(string(op.ID), "target@example.com", urepo, rrepo); err != nil {
		t.Error("user registration failed")
	}
}

func TestUserAddAndDeleteRoleService(t *testing.T) {
	userID := 1
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		suid := string(userID)
		userID++
		return domain.UserID(suid)
	})
	f := domain.NewUserFactory(urepo)
	op := f.Build(domain.UserMail("foo@example.com"))
	tgt := f.Build(domain.UserMail("bar@example.com"))

	rrepo := &domain.RoleRepository{}

	opID := string(op.ID)
	tgtID := string(tgt.ID)

	if err := auth.UserAddRoleService(opID, tgtID, []string{}, urepo, rrepo); err == nil {
		t.Error("operator not found")
	}
	urepo.Save(op)

	if err := auth.UserAddRoleService(opID, tgtID, []string{}, urepo, rrepo); err == nil {
		t.Error("target not found")
	}
	urepo.Save(tgt)

	if err := auth.UserAddRoleService(opID, tgtID, []string{}, urepo, rrepo); err == nil {
		t.Error("role list required")
	}

	if err := auth.UserAddRoleService(opID, tgtID, []string{"!!"}, urepo, rrepo); err == nil {
		t.Error("invalid role exists")
	}

	roleEditor := strconv.Itoa(int(domain.RoleEditor))

	if err := auth.UserAddRoleService(opID, tgtID, []string{roleEditor}, urepo, rrepo); err == nil {
		t.Error("operator has no permission")
	}

	op = op.AddRole(domain.RoleManageUser)
	urepo.Save(op)

	if err := auth.UserAddRoleService(opID, tgtID, []string{roleEditor}, urepo, rrepo); err != nil {
		t.Error("this operation should be valid")
	}

	tgt = urepo.Find(string(tgt.ID))
	if len(tgt.Roles) != 2 {
		t.Error("Editor role should be added")
	}
	if tgt.Roles[1] != domain.RoleEditor {
		t.Error("Editor role should be added")
	}

	if err := auth.UserDeleteRoleService(opID, tgtID, []string{roleEditor}, urepo, rrepo); err != nil {
		t.Error("this operation should be valid")
	}

	tgt = urepo.Find(string(tgt.ID))
	if len(tgt.Roles) != 1 {
		t.Error("Editor role should be deleted")
	}
	if tgt.Roles[0] != domain.RoleViewer {
		t.Error("Editor role should be deleted")
	}

}

func TestSelfRoleEdit(t *testing.T) {
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID("1")
	})
	rrepo := &domain.RoleRepository{}

	f := domain.NewUserFactory(urepo)
	op := f.Build(domain.UserMail("foo@example.com"))
	op = op.AddRole(domain.RoleManageUser)
	urepo.Save(op)

	opID := string(op.ID)
	roleEditor := strconv.Itoa(int(domain.RoleEditor))

	if err := auth.UserAddRoleService(opID, opID, []string{roleEditor}, urepo, rrepo); err != nil {
		t.Error("this operation should be valid")
	}

	op = urepo.Find(opID)
	if len(op.Roles) != 3 {
		t.Error("Viewer, Editor, Manager has attached")
	}
	if op.Roles[2] != domain.RoleEditor {
		t.Error("latest attached role is Editor")
	}

	if err := auth.UserDeleteRoleService(opID, opID, []string{roleEditor}, urepo, rrepo); err != nil {
		t.Error("this operation should be valid")
	}

	op = urepo.Find(opID)
	if len(op.Roles) != 2 {
		t.Error("Viewer, Manager has attached")
	}

	if op.Roles[0] != domain.RoleViewer {
		t.Error("first attached is Viewer")
	}
	if op.Roles[1] != domain.RoleManageUser {
		t.Error("second attached is Manager")
	}
}

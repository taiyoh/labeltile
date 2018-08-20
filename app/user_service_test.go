package app_test

import (
	"strconv"
	"testing"

	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestUserRegisterService(t *testing.T) {
	c := mock.LoadContainerImpl()
	f := domain.NewUserFactory(c.UserRepository())
	op := f.Build(domain.UserMail("foo@example.com"))
	c.UserRepository().Save(op)

	if err := app.UserRegisterService(string(op.ID), "target@example.com", c); err == nil {
		t.Error("operator has no permission")
	}

	op = op.AddRole(domain.RoleManageUser)
	c.UserRepository().Save(op)

	if err := app.UserRegisterService(string(op.ID), "foo@example.com", c); err == nil {
		t.Error("already registered")
	}

	if err := app.UserRegisterService(string(op.ID), "target@example.com", c); err != nil {
		t.Error("user registration failed")
	}
}

func TestUserFindService(t *testing.T) {
	c := mock.LoadContainerImpl()
	f := domain.NewUserFactory(c.UserRepository())
	op := f.Build(domain.UserMail("foo@example.com"))
	op = op.AddRole(domain.RoleManageUser)
	c.UserRepository().Save(op)

	if u := app.UserFindService(string(op.ID+"foo"), c); u != nil {
		t.Error("wrong ID")
	}

	if u := app.UserFindService(string(op.ID), c); u == nil {
		t.Error("finding user failed")
	} else {
		if u.ID != string(op.ID) {
			t.Error("wrong user returns")
		}
		if len(u.Roles) != 2 {
			t.Error("registered roles should be 2")
		}
		if u.Roles[0] != "viewer" || u.Roles[1] != "manager" {
			t.Error("wrong roles registered")
		}
	}
}

func TestUserAddAndDeleteRoleService(t *testing.T) {
	c := mock.LoadContainerImpl()
	f := domain.NewUserFactory(c.UserRepository())
	op := f.Build(domain.UserMail("foo@example.com"))
	tgt := f.Build(domain.UserMail("bar@example.com"))

	opID := string(op.ID)
	tgtID := string(tgt.ID)

	if err := app.UserAddRoleService(opID, tgtID, []string{}, c); err == nil {
		t.Error("operator not found")
	}
	c.UserRepository().Save(op)

	if err := app.UserAddRoleService(opID, tgtID, []string{}, c); err == nil {
		t.Error("target not found")
	}
	c.UserRepository().Save(tgt)

	if err := app.UserAddRoleService(opID, tgtID, []string{}, c); err == nil {
		t.Error("role list required")
	}

	if err := app.UserAddRoleService(opID, tgtID, []string{"!!"}, c); err == nil {
		t.Error("invalid role exists")
	}

	roleEditor := strconv.Itoa(int(domain.RoleEditor))

	if err := app.UserAddRoleService(opID, tgtID, []string{roleEditor}, c); err == nil {
		t.Error("operator has no permission")
	}

	op = op.AddRole(domain.RoleManageUser)
	c.UserRepository().Save(op)

	if err := app.UserAddRoleService(opID, tgtID, []string{roleEditor}, c); err != nil {
		t.Error("this operation should be valid")
	}

	tgt = c.UserRepository().Find(string(tgt.ID))
	if len(tgt.Roles) != 2 {
		t.Error("Editor role should be added")
	}
	if tgt.Roles[1] != domain.RoleEditor {
		t.Error("Editor role should be added")
	}

	if err := app.UserDeleteRoleService(opID, tgtID, []string{roleEditor}, c); err != nil {
		t.Error("this operation should be valid")
	}

	tgt = c.UserRepository().Find(string(tgt.ID))
	if len(tgt.Roles) != 1 {
		t.Error("Editor role should be deleted")
	}
	if tgt.Roles[0] != domain.RoleViewer {
		t.Error("Editor role should be deleted")
	}
}

func TestSelfRoleEdit(t *testing.T) {
	c := mock.LoadContainerImpl()

	f := domain.NewUserFactory(c.UserRepository())
	op := f.Build(domain.UserMail("foo@example.com"))
	op = op.AddRole(domain.RoleManageUser)
	c.UserRepository().Save(op)

	opID := string(op.ID)
	roleEditor := strconv.Itoa(int(domain.RoleEditor))

	if err := app.UserAddRoleService(opID, opID, []string{roleEditor}, c); err != nil {
		t.Error("this operation should be valid")
	}

	op = c.UserRepository().Find(opID)
	if len(op.Roles) != 3 {
		t.Error("Viewer, Editor, Manager has attached")
	}
	if op.Roles[2] != domain.RoleEditor {
		t.Error("latest attached role is Editor")
	}

	if err := app.UserDeleteRoleService(opID, opID, []string{roleEditor}, c); err != nil {
		t.Error("this operation should be valid")
	}

	op = c.UserRepository().Find(opID)
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

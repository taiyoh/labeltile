package domain

import "testing"

func TestUser(t *testing.T) {
	u := NewUser(UserID("1"), UserRoleID("viewer"))
	if u.ID != UserID("1") {
		t.Error("user id should be 1")
	}
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != UserRoleID("viewer") {
		t.Error("user role should be viewer")
	}

	u.AddRole(UserRoleID("editor"))
	if len(u.Roles) != 2 {
		t.Error("user roles count should be 2")
	}
	if u.Roles[1] != UserRoleID("editor") {
		t.Error("user role should be editor")
	}

	u.DeleteRole(UserRoleID("viewer"))
	if len(u.Roles) != 1 {
		t.Error("user roles count should be 1")
	}
	if u.Roles[0] != UserRoleID("editor") {
		t.Error("user role should be only editor")
	}

}

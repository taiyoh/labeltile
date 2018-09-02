package infra_test

import (
	"errors"
	"testing"

	"github.com/taiyoh/labeltile/app"

	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestNewUserRepository(t *testing.T) {
	db := mock.LoadDatabase()
	repo := infra.NewUserRepository(db)
	db.MutateResult(1, nil, 1, nil)
	if id := repo.DispenseID(); id != domain.UserID("1") {
		t.Error("invalid id returns")
	}
	captures := db.Captures()
	if len(captures) != 3 {
		t.Error("captured SQL is not 3:", len(captures))
	}
	for idx, q := range []string{
		"BEGIN",
		"INSERT INTO user_id_dispenser (id) VALUES (default)",
		"COMMIT",
	} {
		if captures[idx].Query != q {
			t.Errorf("idx[%d] wrong captured: %s", idx, captures[idx].Query)
		}
	}
	db.ClearCapture()

	db.MutateError(errors.New("hoge"))
	if id := repo.DispenseID(); id != "" {
		t.Error("valid id returns")
	}
	captures = db.Captures()
	if len(captures) != 3 {
		t.Error("captured SQL is not 3:", len(captures))
	}
	for idx, q := range []string{
		"BEGIN",
		"INSERT INTO user_id_dispenser (id) VALUES (default)",
		"ROLLBACK",
	} {
		if captures[idx].Query != q {
			t.Errorf("idx[%d] wrong captured: %s", idx, captures[idx].Query)
		}
	}
	db.ClearCapture()

	db.MutateResult(1, errors.New("fuga"), 1, nil)
	if id := repo.DispenseID(); id != "" {
		t.Error("valid id returns")
	}
}

func TestUserRepository__Find(t *testing.T) {
	db := mock.LoadDatabase()
	repo := infra.NewUserRepository(db)

	if user := repo.Find("1"); user != nil {
		t.Error("unknown user found")
	}

	db.Add("user", app.DatabaseSelectResult{
		"id":   "1",
		"mail": "foo@example.com",
	})

	user := repo.Find("1")
	if user == nil {
		t.Error("user not found")
	}
	if user.ID != domain.UserID("1") {
		t.Error("wrong user ID")
	}
	if user.Mail != domain.UserMail("foo@example.com") {
		t.Error("wrong mail address")
	}
	if len(user.Roles) != 0 {
		t.Error("unknown roles attached")
	}

	db.Add(
		"user_role",
		app.DatabaseSelectResult{"user_id": "1", "role_id": "1"},
		app.DatabaseSelectResult{"user_id": "1", "role_id": "2"},
		app.DatabaseSelectResult{"user_id": "2", "role_id": "3"},
	)

	db.ClearCapture()

	user = repo.Find("1")
	if user.ID != domain.UserID("1") {
		t.Error("wrong user ID")
	}
	if user.Mail != domain.UserMail("foo@example.com") {
		t.Error("wrong mail address")
	}
	if len(user.Roles) != 2 {
		t.Error("attached role count is wrong")
	}

	for idx, r := range []domain.RoleID{domain.RoleID(1), domain.RoleID(2)} {
		if user.Roles[idx] != r {
			t.Errorf("idx[%d] wrong role attached: %#v\n", idx, user.Roles[idx])
		}
	}

	caps := db.Captures()
	if len(caps) != 2 {
		t.Error("SQL call count is wrong")
	}
	for idx, q := range []string{
		"SELECT * FROM user WHERE id = ? LIMIT 1",
		"SELECT * FROM user_role WHERE user_id IN (?) ORDER BY id ASC",
	} {
		if caps[idx].Query != q {
			t.Errorf("idx[%d] wrong query: %s\n", idx, caps[idx].Query)
		}
	}
}

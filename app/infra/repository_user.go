package infra

import (
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
)

// UserRepository provides implementation for User domain model's data load and store
type UserRepository struct {
	domain.UserRepository
	db app.Database
}

// NewUserRepository returns implemented UserRepository object
func NewUserRepository(db app.Database) *UserRepository {
	return &UserRepository{db: db}
}

// DispenseID returns new ID for User domain model
func (u *UserRepository) DispenseID() domain.UserID {
	id := dispenseID(u.db, "user")
	if id == "" {
		return ""
	}
	return domain.UserID(id)
}

func (u *UserRepository) findRolesByIDs(ids []string) map[string][]domain.RoleID {
	rolesByID := map[string][]domain.RoleID{}
	roles, _ := u.db.Select(sq.Select("*").From("user_role").Where(sq.Eq{"user_id": ids}).OrderBy("id ASC").ToSql())
	for _, role := range roles {
		userID := role["user_id"]
		roleID, _ := strconv.Atoi(role["role_id"])
		rolesID, exists := rolesByID[userID]
		if !exists {
			rolesID = []domain.RoleID{}
		}
		rolesID = append(rolesID, domain.RoleID(roleID))
		rolesByID[userID] = rolesID
	}
	return rolesByID
}

// Find returns User domain model object selecting by id
func (u *UserRepository) Find(id string) *domain.User {
	userRows, _ := u.db.Select(sq.Select("*").From("user").Where(sq.Eq{"id": id}).Limit(1).ToSql())
	if userRows == nil || len(userRows) == 0 {
		return nil
	}
	userRow := userRows[0]

	rolesByID := u.findRolesByIDs([]string{userRow["id"]})
	roles, exists := rolesByID[id]
	if !exists {
		roles = []domain.RoleID{}
	}
	return &domain.User{
		ID:    domain.UserID(userRow["id"]),
		Mail:  domain.UserMail(userRow["mail"]),
		Roles: roles,
	}
}

func (u *UserRepository) exists(id string) bool {
	results, err := u.db.Select(sq.Select("COUNT(id) as user_count").From("user").Where(sq.Eq{"id": id}).ToSql())
	if err != nil || results == nil {
		return false
	}
	cnt, _ := strconv.Atoi(results[0]["user_count"])
	return cnt > 0
}

func (u *UserRepository) Save(du *domain.User) {
	txn := u.db.NewTransaction()
	defer txn.Rollback()
	userID := string(du.ID)
	if !u.exists(userID) {
		uq := sq.Select("*").From("user").Where(sq.Eq{"id": userID}).Limit(1).Suffix("FOR UPDATE")
		if _, err := txn.Select(uq.ToSql()); err != nil {
			return
		}
	}

	iq := sq.Insert("user").Columns("id", "mail").Values(userID, string(du.Mail)).Suffix("ON DUPLICATE KEY UPDATE mail = VALUES(mail)")
	if _, err := txn.Mutate(iq.ToSql()); err != nil {
		return
	}

	txn.Mutate(sq.Delete("user_role").Where(sq.Eq{"user_id": userID}).ToSql())

	rq := sq.Insert("user_role").Columns("user_id", "role_id")
	for _, r := range du.Roles {
		rq = rq.Values(userID, string(r))
	}
	txn.Mutate(rq.ToSql())

	txn.Commit()
}

func (u *UserRepository) FindByMail(mailAddress string) *domain.User {
	userRows, _ := u.db.Select(sq.Select("*").From("user").Where(sq.Eq{"mail": mailAddress}).Limit(1).ToSql())
	if userRows == nil || len(userRows) == 0 {
		return nil
	}
	userRow := userRows[0]

	rolesByID := u.findRolesByIDs([]string{userRow["id"]})
	roles, exists := rolesByID[userRow["id"]]
	if !exists {
		roles = []domain.RoleID{}
	}
	return &domain.User{
		ID:    domain.UserID(userRow["id"]),
		Mail:  domain.UserMail(userRow["mail"]),
		Roles: roles,
	}
}

func (u *UserRepository) FindMulti(ids []string) []*domain.User {
	users := []*domain.User{}
	userRows, _ := u.db.Select(sq.Select("*").From("user").Where(sq.Eq{"id": ids}).ToSql())
	if userRows == nil || len(userRows) == 0 {
		return users
	}

	rolesByID := u.findRolesByIDs(ids)
	for _, user := range userRows {
		uid := user["id"]
		roles, exists := rolesByID[uid]
		if !exists {
			roles = []domain.RoleID{}
		}
		du := &domain.User{
			ID:    domain.UserID(uid),
			Mail:  domain.UserMail(user["mail"]),
			Roles: roles,
		}
		users = append(users, du)
	}
	return users
}

package app

import (
	"errors"

	"github.com/taiyoh/labeltile/app/domain"
)

// UserRegisterService provides user registration application service
func UserRegisterService(opid, mail string, container interface {
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}) error {
	uspec := domain.NewUserSpecification(container.UserRepository())
	rspec := domain.NewRoleSpecification(container.RoleRepository())
	urepo := container.UserRepository()
	op := urepo.Find(opid)
	if op == nil {
		return errors.New("operator not found")
	}
	if err := rspec.SpecifyRegisterUser(op); err != nil {
		return err
	}
	if err := uspec.SpecifyUserRegistration(mail); err != nil {
		return err
	}
	factory := domain.NewUserFactory(urepo)
	u := factory.Build(domain.UserMail(mail))
	urepo.Save(u)

	return nil
}

func loadOperatorAndTarget(opid, tgtid string, urepo domain.UserRepository) (*domain.User, *domain.User, error) {
	var op, tgt *domain.User
	for _, u := range urepo.FindMulti([]string{opid, tgtid}) {
		if u.ID == domain.UserID(opid) {
			op = u
		}
		if u.ID == domain.UserID(tgtid) {
			tgt = u
		}
	}
	if op == nil {
		return nil, nil, errors.New("operator not found")
	}
	if tgt == nil {
		return nil, nil, errors.New("target not found")
	}

	return op, tgt, nil
}

// UserAddRoleService provides attaching role to user
func UserAddRoleService(opid, tgtid string, roles []string, container interface {
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}) error {
	var op, tgt *domain.User
	var roleIDs []domain.RoleID
	var err error

	op, tgt, err = loadOperatorAndTarget(opid, tgtid, container.UserRepository())
	if err != nil {
		return err
	}

	spec := domain.NewRoleSpecification(container.RoleRepository())

	roleIDs, err = spec.ConvertRoleToID(roles)
	if err != nil {
		return err
	}
	if err := spec.SpecifyAddRole(op, tgt, roleIDs); err != nil {
		return err
	}

	for _, rid := range roleIDs {
		tgt = tgt.AddRole(rid)
	}
	container.UserRepository().Save(tgt)

	return nil
}

// UserDeleteRoleService provides detaching role from user
func UserDeleteRoleService(opid, tgtid string, roles []string, container interface {
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}) error {
	var op, tgt *domain.User
	var roleIDs []domain.RoleID
	var err error

	op, tgt, err = loadOperatorAndTarget(opid, tgtid, container.UserRepository())
	if err != nil {
		return err
	}

	spec := domain.NewRoleSpecification(container.RoleRepository())

	roleIDs, err = spec.ConvertRoleToID(roles)
	if err != nil {
		return err
	}
	if err := spec.SpecifyDeleteRole(op, tgt, roleIDs); err != nil {
		return err
	}

	for _, rid := range roleIDs {
		tgt = tgt.DeleteRole(rid)
	}
	container.UserRepository().Save(tgt)

	return nil
}

// UserFindService provides retrieving user data using given id
func UserFindService(id string, container interface {
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}) *UserDTO {
	u := container.UserRepository().Find(id)
	if u == nil {
		return nil
	}
	roles := []string{}
	for _, role := range container.RoleRepository().FindAll(u.Roles) {
		roles = append(roles, role.Name)
	}
	return &UserDTO{
		ID:    string(u.ID),
		Mail:  string(u.Mail),
		Roles: roles,
	}
}

// UserAuthorizeService provides authorized user registration by google
func UserAuthorizeService(code string, container interface {
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
	OAuth2Google() OAuth2Google
}) (*UserDTO, error) {
	info, err := container.OAuth2Google().GetTokenInfo(code)
	if err != nil {
		return nil, err
	}
	repo := container.UserRepository()
	email := info.Email()
	user := repo.FindByMail(email)
	if user == nil {
		user = domain.NewUserFactory(repo).Build(domain.UserMail(email))
	}
	repo.Save(user)
	roles := []string{}
	for _, role := range container.RoleRepository().FindAll(user.Roles) {
		roles = append(roles, role.Name)
	}
	resUser := &UserDTO{
		ID:    string(user.ID),
		Mail:  string(user.Mail),
		Roles: roles,
	}
	return resUser, nil
}

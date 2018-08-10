package auth

import (
	"errors"

	"github.com/taiyoh/labeltile/auth/domain"
)

// UserRegisterService provides user registration application service
func UserRegisterService(opid, mail string, urepo domain.UserRepository, rrepo *domain.RoleRepository) error {
	uspec := domain.NewUserSpecification(urepo)
	rspec := domain.NewRoleSpecification(rrepo)
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
func UserAddRoleService(opid, tgtid string, roles []string, urepo domain.UserRepository, rrepo *domain.RoleRepository) error {
	var op, tgt *domain.User
	var roleIDs []domain.RoleID
	var err error

	op, tgt, err = loadOperatorAndTarget(opid, tgtid, urepo)
	if err != nil {
		return err
	}

	spec := domain.NewRoleSpecification(rrepo)

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
	urepo.Save(tgt)

	return nil
}

// UserDeleteRoleService provides detaching role from user
func UserDeleteRoleService(opid, tgtid string, roles []string, urepo domain.UserRepository, rrepo *domain.RoleRepository) error {
	var op, tgt *domain.User
	var roleIDs []domain.RoleID
	var err error

	op, tgt, err = loadOperatorAndTarget(opid, tgtid, urepo)
	if err != nil {
		return err
	}

	spec := domain.NewRoleSpecification(rrepo)

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
	urepo.Save(tgt)

	return nil
}

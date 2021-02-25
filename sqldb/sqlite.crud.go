package sqldb

import (
	"fmt"
	"local/auth-svc/model"
)

func UpdateUser(user model.User) (model.User, error) {
	database := SQLConnect()

	err := UpdateAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func DeleteUser(user model.User) error {
	database := SQLConnect()

	err := DeleteAccount(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func ListAllAccounts() ([]model.User, error) {
	database := SQLConnect()

	users, err := GetAllAccounts(database)
	if err != nil {
		fmt.Print(err)
	}
	return users, err
}

func AddNewAccount(user model.User) (model.User, error) {
	database := SQLConnect()

	err := AddAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func AddNewRole(role model.Role) (model.Role, error) {
	database := SQLConnect()

	err := AddRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return role, err
}

func ListAllRoles() ([]model.Role, error) {
	database := SQLConnect()

	roles, err := GetAllRoles(database)
	if err != nil {
		fmt.Print(err)
	}
	return roles, err
}

func DeleteRole(role model.Role) error {
	database := SQLConnect()

	err := RemoveRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func UpdateRole(role model.Role) (model.Role, error) {
	database := SQLConnect()

	err := UpdateCurrentRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return role, err
}

package sqldb

import (
	"fmt"
	"local/auth-svc/model"
)

func UpdateUser(user model.User) (model.User, error) {
	database := Connect2Mysql(dbName)

	err := UpdateAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func SetAccountState(user model.User) (model.User, error) {
	database := Connect2Mysql(dbName)

	err := toggleAccountStatus(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func UpdateAccountPassword(user model.User) (model.User, error) {
	database := Connect2Mysql(dbName)

	err := UpdatePassword(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func DeleteUser(user model.User) error {
	database := Connect2Mysql(dbName)

	err := DeleteAccount(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func ListAllAccounts() ([]model.User, error) {
	database := Connect2Mysql(dbName)

	users, err := GetAllAccounts(database)
	if err != nil {
		fmt.Print(err)
	}
	return users, err
}

func AddNewAccount(user model.User) (model.User, error) {
	database := Connect2Mysql(dbName)

	err := AddAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func AddNewRole(role model.Role) (model.Role, error) {
	database := Connect2Mysql(dbName)

	err := AddRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return role, err
}

func ListAllRoles() ([]model.Role, error) {
	database := Connect2Mysql(dbName)

	roles, err := GetAllRoles(database)
	if err != nil {
		fmt.Print(err)
	}
	return roles, err
}

func DeleteRole(role model.Role) error {
	database := Connect2Mysql(dbName)

	err := RemoveRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func UpdateRole(role model.Role) (model.Role, error) {
	database := Connect2Mysql(dbName)

	err := UpdateCurrentRole(database, role)
	if err != nil {
		fmt.Print(err)
	}
	return role, err
}

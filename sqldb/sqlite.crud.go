package sqldb

import (
	"fmt"
	"local/auth-svc/handler"
)

func FindUserByUserName(user handler.User) (handler.User, error) {
	database := SQLConnect()

	query := "SELECT id, username,display_name,email,role,password FROM accounts WHERE username = ?"
	user, err := QueryByParam(database, query, user.Username)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func UpdateUser(user handler.User) (handler.User, error) {
	database := SQLConnect()

	err := UpdateAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func DeleteUser(user handler.User) error {
	database := SQLConnect()

	err := DeleteAccount(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

func ListAllAccounts() ([]handler.User, error) {
	database := SQLConnect()

	users, err := GetAllAccounts(database)
	if err != nil {
		fmt.Print(err)
	}
	return users, err
}

func AddNewAccount(user handler.User) (handler.User, error) {
	database := SQLConnect()

	err := AddAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

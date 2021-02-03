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

func UpdateUserByUserId(user handler.User) (handler.User, error) {
	database := SQLConnect()

	err := UpdateAccountInfo(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func DeleteUserByUserId(user handler.User) error {
	database := SQLConnect()

	err := DeleteAccount(database, user)
	if err != nil {
		fmt.Print(err)
	}
	return err
}

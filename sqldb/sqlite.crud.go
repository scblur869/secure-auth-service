package sqldb

import (
	"fmt"
	"local/auth-svc/handler"
)

func FindUserByUserName(userName string) (handler.User, error) {
	database := SQLConnect()

	query := "SELECT id, username,display_name,email,role,password WHERE username = ?"
	user, err := QueryByParam(database, query, userName)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

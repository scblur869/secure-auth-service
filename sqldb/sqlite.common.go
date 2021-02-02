package sqldb

import (
	"database/sql"
	"fmt"
	"local/auth-svc/handler"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbpath = "data/db"
const db = dbpath + "/auth.db"

func SQLConnect() *sql.DB {
	database, err := sql.Open("sqlite3", db)
	if err != nil {
		fmt.Println(err)
	}
	return database
}

func InitializeDatabase() {
	_, err := os.Stat(dbpath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dbpath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	database, err := sql.Open("sqlite3", db)
	if err != nil {
		fmt.Println(err)
	}
	statement, err :=
		database.Prepare("CREATE TABLE IF NOT EXISTS accounts (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, display_name TEXT, email TEXT, role TEXT, password TEXT, created TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec()

	database.Close()
}

func QueryByParam(database *sql.DB, query string, param string) (handler.User, error) {

	var rowData handler.User
	rows, err := database.Query(query, param)
	var id int
	var username string
	var display_name string
	var email string
	var role string
	var password string

	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &username, &display_name, &email, &role, &password)
		if err != nil {
			fmt.Println(err)
		}

		rowData.ID = id
		rowData.Username = username
		rowData.DisplayName = display_name
		rowData.Email = email
		rowData.Role = role
		rowData.Password = password

	}
	rows.Close()
	database.Close()
	return rowData, err
}

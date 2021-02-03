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

func UpdateAccountInfo(db *sql.DB, user handler.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE accounts SET username=?,display_name=?,email=?,role=?,password=? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Username, user.DisplayName, user.Email, user.Role, user.Password, user.ID); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func DeleteAccount(db *sql.DB, user handler.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("DELETE FROM accounts WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.ID)
	if err != nil {
		fmt.Println(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(affect)
	db.Close()
	return err
}

func AddAccountInfo(db *sql.DB, user handler.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO accounts(username,display_name,email,role,password) VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Username, user.DisplayName, user.Email, user.Role, user.Password); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func GetAllAccounts(database *sql.DB) ([]handler.User, error) {

	var rowData handler.User
	var results []handler.User
	query := "SELECT id, username, display_name, email, role FROM accounts"
	rows, err := database.Query(query)
	var id int
	var username string
	var display_name string
	var email string
	var role string

	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &username, &display_name, &email, &role)
		if err != nil {
			fmt.Println(err)
		}

		rowData.ID = id
		rowData.Username = username
		rowData.DisplayName = display_name
		rowData.Email = email
		rowData.Role = role
		results = append(results, rowData)

	}
	rows.Close()
	database.Close()
	return results, err
}

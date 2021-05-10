package sqldb

import (
	"database/sql"
	"fmt"
	"local/auth-svc/model"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "authdb"

func InitializeDatabase() {

	database := Connect2Mysql(dbName)

	// creates the accounts table
	statement, err :=
		database.Prepare("CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY AUTO_INCREMENT, username VARCHAR(255) UNIQUE NOT NULL, display_name TEXT NOT NULL, email TEXT NOT NULL, role TEXT NOT NULL, is_enabled INTEGER DEFAULT 0, password TEXT NOT NULL, created datetime default CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Println(err)
	}
	acct, err := statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(acct.RowsAffected())

	//creates the initial admin user
	insert, err :=
		database.Prepare("INSERT IGNORE INTO accounts (username, display_name, email, role, is_enabled, password) VALUES(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	ins, err := insert.Exec("admin", "IDP admin user account", "admin@local", "admin", 1, os.Getenv("ADMIN_PASS"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ins.RowsAffected())
	fmt.Println("admin account created")
	//creates the roles table
	roleStmt, err :=
		database.Prepare("CREATE TABLE IF NOT EXISTS roles (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255) UNIQUE NOT NULL, display_name TEXT NOT NULL, description TEXT NOT NULL, created datetime default CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Println(err)
	}
	rol, err := roleStmt.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rol.RowsAffected())

	role_ins, err :=
		database.Prepare("INSERT IGNORE INTO roles (name, display_name, description) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	role_res, err := role_ins.Exec("admin", "Admin Role", "IDP Admin Role")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(role_res.RowsAffected())

	grole_ins, err :=
		database.Prepare("INSERT IGNORE INTO roles (name, display_name, description) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	grole_res, err := grole_ins.Exec("guest", "Guest Role", "IDP Guest Role")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(grole_res.RowsAffected())

	database.Close()
}

func QueryByParam(database *sql.DB, query string, param string) (model.User, error) {
	var rowData model.User
	var id int
	var username string
	var display_name string
	var email string
	var role string
	var is_enabled int // 0 = false / 1 = true
	var password string
	err := database.QueryRow(query, param).Scan(&id, &username, &display_name, &email, &role, &is_enabled, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return rowData, err
		} else {
			fmt.Println(err)
		}
	}

	rowData.ID = id
	rowData.Username = username
	rowData.DisplayName = display_name
	rowData.Email = email
	rowData.Role = role
	rowData.IsEnabled = is_enabled // 0 = false / 1 = true
	rowData.Password = password

	database.Close()
	return rowData, err
}

func UpdateAccountInfo(db *sql.DB, user model.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE accounts SET username=?,display_name=?,email=?,role=? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Username, user.DisplayName, user.Email, user.Role, user.ID); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func UpdatePassword(db *sql.DB, user model.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE accounts SET password=? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Password, user.ID); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func toggleAccountStatus(db *sql.DB, user model.User) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE accounts SET is_enabled=? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.IsEnabled, user.ID); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func DeleteAccount(db *sql.DB, user model.User) error {
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

func AddAccountInfo(db *sql.DB, user model.User) error {
	user.IsEnabled = 0         //accounts are disabled by default
	user.Password = "password" //default password
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO accounts(username,display_name,email,role,is_enabled,password) VALUES(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(user.Username, user.DisplayName, user.Email, user.Role, user.IsEnabled, user.Password); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func GetAllAccounts(database *sql.DB) ([]model.User, error) {

	var rowData model.User
	var results []model.User
	query := "SELECT id, username, display_name, email, role, is_enabled, password FROM accounts"
	rows, err := database.Query(query)
	var id int
	var username string
	var display_name string
	var email string
	var role string
	var is_enabled int
	var password string
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &username, &display_name, &email, &role, &is_enabled, &password)
		if err != nil {
			fmt.Println(err)
		}

		rowData.ID = id
		rowData.Username = username
		rowData.DisplayName = display_name
		rowData.Email = email
		rowData.Role = role
		rowData.IsEnabled = is_enabled
		rowData.Password = password
		results = append(results, rowData)

	}
	rows.Close()
	database.Close()
	return results, err
}

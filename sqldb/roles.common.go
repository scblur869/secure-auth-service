package sqldb

import (
	"database/sql"
	"fmt"
	"local/auth-svc/model"
	"log"
)

func RoleById(database *sql.DB, query string, param string) (model.Role, error) {

	var rowData model.Role
	var id int
	var name string
	var display_name string
	var description string

	err := database.QueryRow(query, param).Scan(&id, &name, &display_name, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			return rowData, err
		} else {
			log.Fatal(err)
		}
	}

	rowData.ID = id
	rowData.Name = name
	rowData.DisplayName = display_name
	rowData.Description = description
	database.Close()
	return rowData, err
}

func UpdateCurrentRole(db *sql.DB, role model.Role) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE roles SET name=?,display_name=?,description=? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(role.Name, role.DisplayName, role.Description); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func RemoveRole(db *sql.DB, role model.Role) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("DELETE FROM roles WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(role.ID)
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

func AddRole(db *sql.DB, role model.Role) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO roles(name,display_name,description) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(role.Name, role.DisplayName, role.Description); err != nil {
		fmt.Println(err)
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	db.Close()
	return err
}

func GetAllRoles(database *sql.DB) ([]model.Role, error) {

	var rowData model.Role
	var results []model.Role
	query := "SELECT id, name, display_name, description FROM roles"
	rows, err := database.Query(query)
	var id int
	var name string
	var display_name string
	var description string

	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &name, &display_name, &description)
		if err != nil {
			fmt.Println(err)
		}

		rowData.ID = id
		rowData.Name = name
		rowData.DisplayName = display_name
		rowData.Description = description
		results = append(results, rowData)

	}
	rows.Close()
	database.Close()
	return results, err
}

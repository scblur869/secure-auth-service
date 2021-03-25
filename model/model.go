package model

import "database/sql"

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	DisplayName string `json:"displayname"`
	Role        string `json:"role"`
	IsEnabled   int    `json:"isenabled"`
}

type Account struct {
	ID          int            `json:"id"`
	Username    sql.NullString `json:"username"`
	Password    sql.NullString `json:"password"`
	Email       sql.NullString `json:"email"`
	DisplayName sql.NullString `json:"displayname"`
	Role        sql.NullString `json:"role"`
	IsEnabled   sql.NullInt32  `json:"isenabled"`
}
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	Description string `json:"description"`
}

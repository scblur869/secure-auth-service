package model

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	DisplayName string `json:"displayname"`
	Role        string `json:"role"`
	IsEnabled   int    `json:"isenabled"`
}

type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	Description string `json:"description"`
}

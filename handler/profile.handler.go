package handler

import (
	"local/auth-svc/auth"
)

// ProfileHandler struct
type profileHandler struct {
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewProfile(rd auth.AuthInterface, tk auth.TokenInterface) *profileHandler {
	return &profileHandler{rd, tk}
}

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	DisplayName string `json:"displayname"`
	Role        string `json:"role"`
}

//In memory user
var user = User{
	ID:          1,
	Username:    "tsadmin",
	Password:    "tsadmin",
	Email:       "admin@temper-sure.com",
	DisplayName: "Temper-Sure Admin",
	Role:        "admin",
}

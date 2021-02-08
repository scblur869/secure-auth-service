package handler

import (
	"fmt"
	"local/auth-svc/auth"
	"os"

	"github.com/dgrijalva/jwt-go"
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
	Email:       "tsadmin@test.com",
	DisplayName: "Demo Account",
	Role:        "user-role",
}

func resolveClaims(tokenString string) map[string]interface{} {
	c := make(map[string]interface{})
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	c["display_name"] = claims["display_name"]
	c["role"] = claims["role"]
	return c
}

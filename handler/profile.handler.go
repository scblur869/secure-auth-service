package handler

// inspired by https://github.com/victorsteven/jwt-best-practices
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

// The only in-memory user
var user = User{
	ID:          1,
	Username:    "demo",
	Password:    "demo",
	Email:       "demo@test.com",
	DisplayName: "Demo User",
	Role:        "Demo",
}

type Claims struct {
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
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

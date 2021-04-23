package handler

import (
	"encoding/json"
	"fmt"
	crypt "local/auth-svc/_cipher"
	"local/auth-svc/model"
	sql "local/auth-svc/sqldb"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindUserByUserName(user model.User) (model.User, error) {
	database := sql.SQLConnect()

	query := "SELECT id, username, display_name, email, role, is_enabled, password FROM accounts WHERE username = ?"
	user, err := sql.QueryByParam(database, query, user.Username)
	if err != nil {
		fmt.Print(err)
	}
	return user, err
}

func (h *profileHandler) SendLoginCookie(c *gin.Context) {
	var u model.User
	strconv.Itoa(user.ID)
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := FindUserByUserName(u)
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := h.tk.CreateToken(strconv.Itoa(user.ID), user.Email, user.DisplayName, user.Role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := h.rd.CreateAuth(strconv.Itoa(user.ID), ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	jsonString, err := json.Marshal(tokens)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	claims := resolveClaims(ts.AccessToken)
	jsonStr, err := json.Marshal(claims)
	if err != nil {
		fmt.Println(err)
	}
	domain := os.Getenv("COOKIE_DOMAIN")
	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		fmt.Println(err)
	}
	encCookie := crypt.Encrypt(string(jsonString), os.Getenv("AESKEY"))
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("ts-cookie", encCookie, 108000, "", domain, secure, true)
	c.SetCookie("is-logged-in", string(jsonStr), 1800, "", domain, false, false)
	c.JSON(http.StatusOK, "successful")
}

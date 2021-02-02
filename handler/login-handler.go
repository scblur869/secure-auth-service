package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *profileHandler) SendLoginCookie(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	ts, err := h.tk.CreateToken(user.ID, user.Email, user.DisplayName, user.Role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := h.rd.CreateAuth(user.ID, ts)
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
	c.SetCookie("ts-cookie", string(jsonString), 108000, "", "", false, true)
	c.JSON(http.StatusOK, ts.AccessToken)
}

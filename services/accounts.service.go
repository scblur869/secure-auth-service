package service

import (
	"local/auth-svc/handler"
	"local/auth-svc/model"
	"local/auth-svc/sqldb"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAccounts(c *gin.Context) {
	users, err := sqldb.ListAllAccounts()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func AddAccount(c *gin.Context) {

	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := sqldb.AddNewAccount(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func ModifyAccount(c *gin.Context) {
	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := sqldb.UpdateUser(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func ToggleAccount(c *gin.Context) {
	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := sqldb.SetAccountState(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func SetPassword(c *gin.Context) {
	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := sqldb.UpdateAccountPassword(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func RemoveAccount(c *gin.Context) {
	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	err := sqldb.DeleteUser(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "user "+account.Username+" removed")
}

func FindUser(c *gin.Context) {
	var account model.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user, err := handler.FindUserByUserName(account)
	if err != nil {
		c.JSON(http.StatusNoContent, err)

		return
	}
	c.JSON(http.StatusOK, user)
}

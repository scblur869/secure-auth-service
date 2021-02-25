package service

import (
	"local/auth-svc/model"
	"local/auth-svc/sqldb"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListRoles(c *gin.Context) {
	roles, err := sqldb.ListAllRoles()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

func AddRole(c *gin.Context) {

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	role, err := sqldb.AddNewRole(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, role)
}

func ModifyRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	role, err := sqldb.UpdateRole(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, role)
}

func RemoveRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	err := sqldb.DeleteRole(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "role "+role.Name+" removed")
}

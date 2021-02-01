package handler

import (
	"encoding/json"
	"fmt"
	"local/auth-svc/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// logout function
func (h *profileHandler) LogoutSession(c *gin.Context) {
	mapToken := map[string]string{}
	ck, _ := c.Cookie("ts-cookie")

	if err := json.Unmarshal([]byte(ck), &mapToken); err != nil {
		fmt.Println(err)
	}
	token, terr := auth.VerifyTokenMap(mapToken["access_token"])

	if terr != nil {
		fmt.Println("error resolving token map")
	}
	metadata, _ := h.tk.ResolveToken(token)
	if metadata != nil {
		deleteErr := h.rd.DeleteTokens(metadata)
		if deleteErr != nil {
			c.JSON(http.StatusBadRequest, deleteErr.Error())
			return
		}
	}
	c.SetCookie("ts-cookie", "stale", -1, "", "", false, true)
	c.JSON(http.StatusOK, "Successfully logged out")
}

package handler

import (
	"encoding/json"
	"fmt"
	crypt "local/auth-svc/_cipher"
	"local/auth-svc/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// logout function
func (h *profileHandler) LogoutSession(c *gin.Context) {
	mapToken := map[string]string{}
	cookie, ckErr := c.Cookie("ts-cookie")
	if ckErr != nil {
		fmt.Println(ckErr)
		c.JSON(http.StatusBadRequest, "data requirement not met")
		return
	}
	ck, err := crypt.Decrypt(string(cookie), os.Getenv("AESKEY"))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, "data requirement not met")
		return
	}

	if err := json.Unmarshal([]byte(ck), &mapToken); err != nil {
		fmt.Println(err)
	}
	token, terr := auth.VerifyTokenMap(mapToken["access_token"])

	if terr != nil {
		fmt.Println("error resolving token map")
	}
	metadata, merr := h.tk.ResolveToken(token)
	if merr != nil {
		fmt.Println(merr)
		return
	}
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

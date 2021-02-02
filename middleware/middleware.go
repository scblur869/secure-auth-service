package middleware

import (
	"encoding/json"
	"fmt"
	"local/auth-svc/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// deprecated for now, token validation is done in the handler
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		ck, _ := c.Cookie("ts-cookie")

		if err := json.Unmarshal([]byte(ck), &mapToken); err != nil {
			fmt.Println(err)
		}

		terr, _ := auth.VerifyTokenMap(mapToken["access_token"])
		if terr != nil {
			//	fmt.Println(terr)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		fmt.Println("success")
		c.Next()
	}
}

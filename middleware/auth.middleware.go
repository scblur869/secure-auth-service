package middleware

import (
	"encoding/json"
	"fmt"
	"local/auth-svc/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		ck, _ := c.Cookie("ts-cookie")

		if err := json.Unmarshal([]byte(ck), &mapToken); err != nil {
			fmt.Println(err)
		}

		token, terr := auth.VerifyTokenMap(mapToken["access_token"])
		if terr != nil {
			fmt.Println(terr)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		if token.Valid == true {
			fmt.Println("success")
			c.Next()
		}
	}
}

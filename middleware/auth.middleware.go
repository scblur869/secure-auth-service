package middleware

import (
	"encoding/json"
	"fmt"
	crypt "local/auth-svc/_cipher"
	"local/auth-svc/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// maybe a little overkill on the checking and dumping to a 401 ..
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapToken := map[string]string{}
		cookie, err := c.Cookie("ts-cookie")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		ck, err := crypt.Decrypt(string(cookie), os.Getenv("AESKEY"))

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		if err := json.Unmarshal([]byte(ck), &mapToken); err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token, terr := auth.VerifyTokenMap(mapToken["access_token"])
		if terr != nil {
			fmt.Println(terr)
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		if token.Valid == true {
			fmt.Println("token is valid")
			c.Next()
		}
	}
}

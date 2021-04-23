package handler

import (
	"encoding/json"
	"fmt"
	crypt "local/auth-svc/_cipher"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// token refresh function
func (h *profileHandler) RefreshSession(c *gin.Context) {
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
		panic(err)
	}

	refreshToken := mapToken["refresh_token"]
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		userId, roleOk := claims["user_id"].(string)
		userEmail, _ := claims["email"].(string)
		displayName, _ := claims["display_name"].(string)
		userRole, _ := claims["role"].(string)
		if !roleOk {
			fmt.Println("roleOk ", roleOk)
			c.JSON(http.StatusUnprocessableEntity, "unauthorized")
			return
		}
		//Delete the previous Refresh Token
		delErr := h.rd.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := h.tk.CreateToken(userId, userEmail, displayName, userRole)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := h.rd.CreateAuth(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
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
		c.JSON(http.StatusCreated, "successful")
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

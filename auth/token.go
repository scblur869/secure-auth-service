package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type tokenservice struct{}

func NewToken() *tokenservice {
	return &tokenservice{}
}

type TokenInterface interface {
	CreateToken(userId string, userEmail string, displayName string, role string) (*TokenDetails, error)
	ResolveToken(token *jwt.Token) (*AccessDetails, error)
}

//Token implements the TokenInterface
var _ TokenInterface = &tokenservice{}

func (t *tokenservice) CreateToken(userId string, userEmail string, displayName string, role string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 10).Unix() //expires after 10 min
	td.TokenUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 2).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = userId
	atClaims["email"] = userEmail
	atClaims["display_name"] = displayName
	atClaims["role"] = role
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS384, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["email"] = userEmail
	rtClaims["display_name"] = displayName
	rtClaims["role"] = role
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func VerifyTokenMap(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("possible incorrect algorithm: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//get the token from the request body
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func extract(token *jwt.Token) (*AccessDetails, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		displayname, _ := claims["display_name"].(string)
		role, _ := claims["role"].(string)
		if ok == false || userOk == false {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid:   accessUuid,
				UserId:      userId,
				Email:       email,
				DisplayName: displayname,
				Role:        role,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func (t *tokenservice) ResolveToken(token *jwt.Token) (*AccessDetails, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		email, _ := claims["email"].(string)
		displayname, _ := claims["display_name"].(string)
		role, _ := claims["role"].(string)
		if ok == false || userOk == false {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid:   accessUuid,
				UserId:      userId,
				Email:       email,
				DisplayName: displayname,
				Role:        role,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

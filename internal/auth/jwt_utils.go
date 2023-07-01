package auth

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GenerateToken generates a jwt token for a user
func GenerateToken(userId bson.ObjectId) (string, error) {

	tokenLifespan, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN_HOUR"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId.Hex()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// IsTokenValid validates the token
func IsTokenValid(c *gin.Context) error {
	tokenString := ExtractTokenFromContext(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}
	return nil
}

// ExtractTokenFromContext extracts the token from the request
func ExtractTokenFromContext(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID extracts the token id from the request
func ExtractTokenID(c *gin.Context) (string, error) {

	tokenString := ExtractTokenFromContext(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := claims["userId"].(string)
		return uid, nil
	}
	return "", nil
}

// Logout logs out the user
/*func Logout(c *gin.Context) error {
	_, err := ExtractTokenID(c)
	if err != nil {
		return err
	}

}*/

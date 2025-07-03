package jwt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserIDKey            = "sub"
	AuthorizedKey        = "authorized"
	ExpKey               = "exp"
	JwtExpirationKey     = "JWT_EXPIRES_IN_HOUR"
	RefreshExpirationKey = "JWT_REFRESH_EXPIRES_IN_HOUR"
	JwtSecretKey         = "JWT_SECRET"
	AuthorizationHeader  = "Authorization"
)

type TokenPair struct {
	Token        string             `json:"token"`
	RefreshToken string             `json:"refreshToken"`
	UserID       primitive.ObjectID `json:"-"`
}

// GenerateTokenPair generates a jwt and refresh token
func GenerateTokenPair(userID primitive.ObjectID) (TokenPair, error) {

	tokenLifespanStr := os.Getenv(JwtExpirationKey)
	if tokenLifespanStr == "" {
		return TokenPair{}, fmt.Errorf("JWT_EXPIRES_IN_HOUR environment variable not set")
	}

	tokenLifespan, err := strconv.Atoi(tokenLifespanStr)
	if err != nil {
		return TokenPair{}, fmt.Errorf("invalid JWT_EXPIRES_IN_HOUR value: %v", err)
	}

	if tokenLifespan <= 0 {
		return TokenPair{}, fmt.Errorf("JWT_EXPIRES_IN_HOUR must be a positive number")
	}

	jwtSecret := os.Getenv(JwtSecretKey)
	if jwtSecret == "" {
		return TokenPair{}, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	claims := jwt.MapClaims{}
	claims[AuthorizedKey] = true
	claims[UserIDKey] = userID.Hex()
	claims[ExpKey] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign JWT token: %v", err)
	}

	// Refresh token is a random string
	refreshTokenStr := primitive.NewObjectID().Hex()

	return TokenPair{Token: tokenStr, RefreshToken: refreshTokenStr, UserID: userID}, nil
}

// IsTokenValid validates the token
func IsTokenValid(token string) error {
	jwtSecret := os.Getenv(JwtSecretKey)
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable not set")
	}

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}
	return nil
}

// ExtractTokenFromContext extracts the user id from the bearer token or refresh token
func ExtractTokenFromContext(c *gin.Context) (string, error) {
	token := extractBearerTokenFromContext(c)

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

// extractBearerTokenFromContext extracts the token from the request
func extractBearerTokenFromContext(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get(AuthorizationHeader)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractUserIDFromContext extracts the token id (userID) from the request
func ExtractUserIDFromContext(c *gin.Context) (string, error) {

	tokenString, err := ExtractTokenFromContext(c)
	if err != nil {
		return "", err
	}

	jwtSecret := os.Getenv(JwtSecretKey)
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := claims[UserIDKey].(string)
		return uid, nil
	}
	return "", nil
}

// Logout logs out the user
/*func Logout(c *gin.Context) error {
	_, err := ExtractUserIDFromContext(c)
	if err != nil {
		return err
	}

}*/

package utils

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harshambasta-2001/Steganography_project/internal"
	"golang.org/x/crypto/bcrypt"
)
type JWTClaims struct {
	UserID uint `json:"userId"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user internal.User) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		return "", errors.New("JWT_SECRET environment variable not set")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID: user.ID,
		UserName: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}


func validateToken(tokenString string) (*JWTClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		claims, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}
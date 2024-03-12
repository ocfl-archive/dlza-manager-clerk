package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func tokenValid(c *gin.Context, key string) error {
	tokenString := extractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func extractToken(c *gin.Context) string {
	reqToken := c.Request.Header.Get("Authorization")
	splitStr := strings.Split(reqToken, "Bearer ")
	if len(splitStr) == 2 {
		return splitStr[1]
	}
	return ""
}

func JwtAuthMiddleware(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenValid(c, key)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

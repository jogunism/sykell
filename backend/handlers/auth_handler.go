package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"backend/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("hyunwoo-using-jwt-for-sykell-homeproject")

// AuthMiddleware is a Gin middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrMissingAuthHeader.Error()})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &domain.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		// fmt.Println(claims.Username)

		// claimsJSON, marshalErr := json.MarshalIndent(claims, "", "  ")
		// if marshalErr != nil {
		// 	fmt.Println("Error marshaling claims to JSON:", marshalErr)
		// }
		// fmt.Println(string(claimsJSON))

		if err != nil || !token.Valid {
			fmt.Println("JWT parsing error:", err) // Keep this for detailed error messages
			c.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrTokenInvalid.Error()})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

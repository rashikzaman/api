package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"

	"github.com/gin-gonic/gin"
)

// ClerkMiddleware creates a new Clerk session middleware for Gin
func ClerkMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

		clerk.SetKey(secretKey)

		// Verify the session
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
			//JWKSClient: jwksClient,
		})
		if err != nil {
			fmt.Println(sessionToken)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		usr, err := user.Get(c.Request.Context(), claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("user", usr)

		c.Next()
	}
}

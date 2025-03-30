package middleware

import (
	"fmt"
	"net/http"
	"rashikzaman/api/application"
	"rashikzaman/api/models"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"

	"github.com/gin-gonic/gin"
)

// ClerkMiddleware creates a new Clerk session middleware for Gin
func ClerkMiddleware(app *application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

		clerk.SetKey(app.Config.GetClerkSecretKey())

		// verify the session
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			fmt.Println("middleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		clerkUser, err := user.Get(c.Request.Context(), claims.Subject)
		if err != nil {
			fmt.Println("middleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		user, err := models.GetUserByClerkID(c, app.DB, clerkUser.ID)
		if err != nil {
			fmt.Println("middleware error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

// ClerkMiddlewareForAdmin creates a new Clerk session middleware for Gin specifically for admin
func ClerkMiddlewareForAdmin(app *application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

		clerk.SetKey(app.Config.GetClerkSecretKey())

		// verify the session
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		clerkUser, err := user.Get(c.Request.Context(), claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		user, err := models.GetUserByClerkID(c, app.DB, clerkUser.ID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "this route is protected"})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

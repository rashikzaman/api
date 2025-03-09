package routes

import (
	"net/http"
	"rashikzaman/api/application"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

func PrivateRoutes(r *gin.Engine, app application.Application) {
	routeGroup := r.Group("/private", middleware.ClerkMiddleware(app.Config.GetClerkSecretKey()))

	routeGroup.GET("/test", Test)
}

func Test(c *gin.Context) {
	usr, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
	}

	c.JSON(http.StatusOK, gin.H{"user": usr.(*clerk.User).FirstName})
}

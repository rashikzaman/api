package http

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/routes"

	"github.com/gin-gonic/gin"
)

func RunHTTPServer(app application.Application) {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Static("/uploads", "../uploads")

	// Then register routes
	routes.AuthRoutes(r, app)
	routes.TaskRoutes(r, app)
	routes.CategoryRoutes(r, app)
	routes.SkillsRoutes(r, app)
	routes.AdminRoutes(r, app)
	routes.UserRoutes(r, app)

	_ = r.Run(":" + app.Config.GetHTTPPort())
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

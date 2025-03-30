package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/users", middleware.ClerkMiddleware(&app))

	routeGroup.GET("/me", controller.FetchMe)
	routeGroup.PUT("/me", controller.UpdateMe)
}

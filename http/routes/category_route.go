package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/categories", middleware.ClerkMiddleware(&app))

	routeGroup.GET("/", controller.FetchCategories)
}

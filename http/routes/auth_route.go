package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/auth")

	routeGroup.POST("/register", controller.Register)
	routeGroup.POST("/login", controller.Login)
}

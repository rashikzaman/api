package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SkillsRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/skills", middleware.ClerkMiddleware(&app))

	routeGroup.GET("/", controller.FetchSkills)
}

package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/admin", middleware.ClerkMiddlewareForAdmin(&app))

	routeGroup.GET("/me", controller.GetAdmin)
	routeGroup.GET("/tasks", controller.FetchTasksForAdmin)
	routeGroup.GET("/users", controller.FetchUsersForAdmin)
	routeGroup.PATCH("/users/:id/:action", controller.ApplyActionToUser)
	routeGroup.PATCH("/tasks/:id/:action", controller.ApplyActionToTask)
}

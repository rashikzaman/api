package routes

import (
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"
	middleware "rashikzaman/api/http/middlewares"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/tasks", middleware.ClerkMiddleware(&app))

	routeGroup.GET("/", controller.FetchTasks)
	routeGroup.GET("/me", controller.FetchTasksCreatedByUser)
	routeGroup.GET("/me/subscribed", controller.FetchTasksSubscribedByUser)
	routeGroup.POST("/", controller.CreateTask)
	routeGroup.DELETE("/:id/", controller.DeleteTask)
	routeGroup.GET("/:id", controller.FetchTask)
	routeGroup.PUT("/:id", controller.UpdateTask)
	routeGroup.POST("/:id/apply", controller.ApplyToTask)
	routeGroup.DELETE("/:id/withdraw", controller.WithdrawFromTask)
	routeGroup.GET("/:id/subscribers", controller.GetSubscribersOfTask)
}

package routes

import (
	"fmt"
	"rashikzaman/api/application"
	"rashikzaman/api/http/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, app application.Application) {
	controller := controllers.Controller{
		App: &app,
	}

	routeGroup := r.Group("/webhook")

	routeGroup.POST("/clerk", controller.ClerkWebHook)
	routeGroup.GET("/hello", hello)
}

func hello(c *gin.Context) {
	fmt.Println("hello world")
}

package controllers

import (
	"rashikzaman/api/application"
	"rashikzaman/api/models"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	App *application.Application
}

func GetUser(c *gin.Context) *models.User {
	userValue, _ := c.Get("user")

	user, _ := userValue.(*models.User)

	return user

}

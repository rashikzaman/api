package controllers

import (
	"fmt"
	"net/http"
	"rashikzaman/api/services"

	"github.com/gin-gonic/gin"
)

func (ac *Controller) FetchSkills(c *gin.Context) {
	skills, err := services.FetchSkills(
		c, ac.App.DB)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, skills)
}

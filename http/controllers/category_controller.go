package controllers

import (
	"fmt"
	"net/http"
	"rashikzaman/api/services"

	"github.com/gin-gonic/gin"
)

func (ac *Controller) FetchCategories(c *gin.Context) {
	categories, err := services.FetchCategories(
		c, ac.App.DB)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, categories)
}

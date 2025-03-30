package controllers

import (
	"fmt"
	"net/http"
	"rashikzaman/api/models"
	"rashikzaman/api/services"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func (ac *Controller) FetchMe(c *gin.Context) {
	user := GetUser(c)

	detailedUser, err := services.GetUserByID(c, ac.App.DB, user.ID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, detailedUser)
}

func (ac *Controller) UpdateMe(c *gin.Context) {
	user := GetUser(c)

	userBody := models.User{}

	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err := models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		_, err := services.UpdateMe(c, tx, user, userBody)
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = services.CreateOrUpdateUserLocation(c, tx, &userBody.UserLocations, user.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, userBody)
}

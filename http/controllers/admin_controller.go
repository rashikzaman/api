package controllers

import (
	"fmt"
	"net/http"
	"rashikzaman/api/models"
	"rashikzaman/api/services"
	"rashikzaman/api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ac *Controller) GetAdmin(c *gin.Context) {
	user := GetUser(c)
	c.JSON(http.StatusOK, &user)
}

func (ac *Controller) FetchTasksForAdmin(c *gin.Context) {
	pagination := utils.PaginationConfigFromRequest(c)

	tasks, count, err := services.FetchTasks(
		c, ac.App.DB,
		models.QueryParam{
			Pagination: utils.PaginationConfigFromRequest(c),
			Relations:  []string{"User", "Category", "Media", "SubscribedUsers"},
		},
		services.Filter{},
	)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, struct {
		Count      int         `json:"count"`
		PageNumber int         `json:"pageNumber"`
		Records    interface{} `json:"records"`
	}{
		Count:      count,
		PageNumber: pagination.Page,
		Records:    tasks,
	})
}

func (ac *Controller) FetchUsersForAdmin(c *gin.Context) {
	pagination := utils.PaginationConfigFromRequest(c)

	tasks, count, err := services.FetchUsersForAdmin(
		c, ac.App.DB,
		models.QueryParam{
			Pagination: utils.PaginationConfigFromRequest(c),
		},
	)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, struct {
		Count      int         `json:"count"`
		PageNumber int         `json:"pageNumber"`
		Records    interface{} `json:"records"`
	}{
		Count:      count,
		PageNumber: pagination.Page,
		Records:    tasks,
	})
}

func (ac *Controller) ApplyActionToUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	action := c.Param("action")

	_, err = services.ApplyActionToUser(c, ac.App.DB, userID, action)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusOK)
}

func (ac *Controller) ApplyActionToTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	action := c.Param("action")

	_, err = services.ApplyActionToTask(c, ac.App.DB, taskID, action)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Status(http.StatusOK)
}

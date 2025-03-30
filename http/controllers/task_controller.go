package controllers

import (
	"fmt"
	"net/http"
	"rashikzaman/api/models"
	"rashikzaman/api/services"
	"rashikzaman/api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (ac *Controller) CreateTask(c *gin.Context) {
	task := &models.Task{}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	user := GetUser(c)

	err := models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		err := services.CreateTask(c, tx, task, user.ID)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (ac *Controller) UpdateTask(c *gin.Context) {
	task := &models.Task{}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	user := GetUser(c)

	existingTask, err := services.FetchTaskByID(c, ac.App.DB, taskID, models.QueryParam{})
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if user.ID != existingTask.UserID {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	err = models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		existingTask, err := services.FetchTaskByID(c, tx, taskID, models.QueryParam{})
		if err != nil {
			return err
		}

		_, err = services.UpdateTask(c, tx, existingTask, *task)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (ac *Controller) FetchTasks(c *gin.Context) {
	pagination := utils.PaginationConfigFromRequest(c)
	categoryIDs := c.Query("category_ids")
	skills := c.Query("skills")
	categoryIDArray := []string{}
	skillsArray := []string{}

	if categoryIDs != "" {
		categoryIDArray = strings.Split(categoryIDs, ",")
	}

	if skills != "" {
		skillsArray = strings.Split(skills, ",")
	}

	var (
		longitude, latitude float32
		subscribeUserId     uuid.UUID
	)

	if c.Query("longitude") != "" {
		lng, err := strconv.ParseFloat(c.Query("longitude"), 32)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		longitude = float32(lng)
	}

	if c.Query("latitude") != "" {
		lat, err := strconv.ParseFloat(c.Query("latitude"), 32)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		latitude = float32(lat)
	}

	if c.Query("is_subscribed") == "true" {
		subscribeUserId = GetUser(c).ID
	}

	distance, err := strconv.ParseInt(c.Query("distance"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	tasks, count, err := services.FetchTasks(
		c, ac.App.DB,
		models.QueryParam{
			Pagination: utils.PaginationConfigFromRequest(c),
			Relations:  []string{"User", "Category", "Media", "SubscribedUsers"},
		},
		services.Filter{
			CategoryIDs:      categoryIDArray,
			Skills:           skillsArray,
			Longitude:        longitude,
			Latitude:         latitude,
			SubscribeUserID:  subscribeUserId,
			SearchTerm:       c.Query("search_term"),
			ApplyBlockFilter: true,
			Distance:         int(distance),
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

func (ac *Controller) FetchTasksCreatedByUser(c *gin.Context) {
	user := GetUser(c)
	pagination := utils.PaginationConfigFromRequest(c)

	tasks, count, err := services.FetchTasks(
		c, ac.App.DB,
		models.QueryParam{
			Pagination: utils.PaginationConfigFromRequest(c),
			Relations:  []string{"User", "Category", "Media", "SubscribedUsers"},
		},
		services.Filter{
			CreatedByUserID: user.ID,
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

func (ac *Controller) FetchTasksSubscribedByUser(c *gin.Context) {
	user := GetUser(c)
	pagination := utils.PaginationConfigFromRequest(c)

	tasks, count, err := services.FetchTasks(
		c, ac.App.DB,
		models.QueryParam{
			Pagination: utils.PaginationConfigFromRequest(c),
			Relations:  []string{"User", "Category", "Media", "SubscribedUsers"},
		},
		services.Filter{
			SubscribedByUserID: user.ID,
			ApplyBlockFilter:   true,
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

func (ac *Controller) ApplyToTask(c *gin.Context) {
	user := GetUser(c)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		err := services.ApplyToTask(c, tx, taskID, user.ID)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (ac *Controller) WithdrawFromTask(c *gin.Context) {
	user := GetUser(c)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		err := services.WithdrawFromTask(c, tx, taskID, user.ID)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (ac *Controller) DeleteTask(c *gin.Context) {
	user := GetUser(c)

	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	task, err := services.FetchTaskByID(c, ac.App.DB, taskID, models.QueryParam{})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	//only the user who created the task can delete the task
	if task.UserID != user.ID {
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	err = models.WithTransaction(c, ac.App.DB, func(tx *bun.Tx) error {
		err := services.DeleteTask(c, tx, taskID)
		return err
	})

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (ac *Controller) FetchTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	task, err := services.FetchTaskByID(c, ac.App.DB, taskID, models.QueryParam{
		Relations: []string{"User", "Category", "Media", "SubscribedUsers"},
		Alias:     "task",
	})
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, task)
}

func (ac *Controller) GetSubscribersOfTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userTasks, err := services.FetchSubscribersForTask(c, ac.App.DB, taskID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, userTasks)
}

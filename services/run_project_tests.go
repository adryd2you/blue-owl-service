package services

import (
	"blue-owl-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunProjectTests(c *gin.Context) {
	var runRequest models.TestRequest

	if err := c.BindJSON(&runRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}

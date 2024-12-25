package services

import (
	"blue-owl-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunTests(c *gin.Context) {
	var runRequest models.TestRequest

	if err := c.BindJSON(&runRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if runRequest.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id must not be null"})
		return
	}

	if runRequest.EndpointID != 0 && runRequest.ServiceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id must not be null if endpoint_id is provided"})
		return
	}

	if runRequest.ExampleID != 0 && (runRequest.ServiceID == 0 || runRequest.EndpointID == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id and endpoint_id must not be null if test_id is provided"})
		return
	}

	if runRequest.ServiceID == 0 && runRequest.EndpointID == 0 && runRequest.ExampleID == 0 {
		RunProjectTests(c)
	} else if runRequest.EndpointID == 0 && runRequest.ExampleID == 0 {
		RunServiceTests(c)
	} else if runRequest.ExampleID == 0 {
		RunEndpointTests(c)
	} else {
		err := ExecuteRunSpecificTest(&runRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusCreated, gin.H{
			"message": fmt.Sprintf("Success Run Specific Test by Example ID: %d", runRequest.ExampleID),
		})
	}
}

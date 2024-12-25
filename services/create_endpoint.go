package services

import (
	"blue-owl-service/models"
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAPIEndpoint creates a new API endpoint
func CreateAPIEndpoint(c *gin.Context) {
	// Create a variable to hold the request data
	var createAPIEndpoint models.CreateAPIEndpoint

	// Bind the incoming JSON request to the createAPIEndpoint struct
	if err := c.ShouldBindJSON(&createAPIEndpoint); err != nil {
		// If binding fails, return a 400 Bad Request with an error message
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation (optional): Ensure that essential fields are provided
	if createAPIEndpoint.ServiceID == 0 || createAPIEndpoint.Method == "" || createAPIEndpoint.Name == "" || createAPIEndpoint.Path == "" || createAPIEndpoint.APIType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// var formattedJSON bytes.Buffer
	// // Unmarshal and then Marshal with indentation
	// err := json.Indent(&formattedJSON, []byte(createAPIEndpoint.RequestBody), "", "    ")
	// if err != nil {
	// 	fmt.Println("Error formatting JSON:", err)
	// 	return
	// }
	// createAPIEndpoint.RequestBody = json.RawMessage(formattedJSON.String())

	// Call the repository function to insert the API endpoint into the database
	// (Assuming you have an InsertAPIEndpoint function that handles the insertion)
	createdID, err := repositories.InsertAPIEndpoint(createAPIEndpoint, transport.Db)
	if err != nil {
		// If an error occurs while inserting the API endpoint, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create API endpoint: %v", err)})
		return
	}

	exampleID, err := repositories.InsertAPIExample(createAPIEndpoint, transport.Db, int(createdID))
	if err != nil {
		// If an error occurs while inserting the API endpoint, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create API endpoint: %v", err)})
		return
	}

	responseAIRequest, err := RequestContractAI(createAPIEndpoint)
	if err != nil {
		log.Printf("fail Call AI")
	}
	requestContracts := mapParametersToRequestContracts(responseAIRequest, int(createdID))

	requestID, err := repositories.InsertAPIRequests(requestContracts, transport.Db)
	if err != nil {
		// If an error occurs while inserting the API endpoint, return a 500 Internal Server Error
		log.Printf("fail Insert Request")
	}

	responseAIResponse, err := ResponseContractAI(createAPIEndpoint)
	if err != nil {
		log.Printf("fail Call AI")
	}
	responseContracts := mapParametersToResponseContracts(responseAIResponse, int(createdID))

	responseID, err := repositories.InsertAPIResponse(responseContracts, transport.Db)
	if err != nil {
		// If an error occurs while inserting the API endpoint, return a 500 Internal Server Error
		log.Printf("fail Insert Response")
	}
	// Check if the requestID or responseID is empty or nil, and return a different response if needed
	if len(requestID) == 0 {
		requestID = nil // Or you can set it to an empty array: requestID = []int64{}
	}

	if len(responseID) == 0 {
		responseID = nil // Or you can set it to an empty array: responseID = []int64{}
	}

	// Return a success response with the created API endpoint ID
	c.JSON(http.StatusCreated, gin.H{
		"message":     "API endpoint created successfully",
		"id":          createdID,
		"example_id":  exampleID,
		"request_id":  requestID,
		"response_id": responseID,
	})
}

func mapParametersToRequestContracts(parameters []Parameter, endpointID int) []models.RequestContract {
	var contracts []models.RequestContract

	for i, param := range parameters {
		// Map each Parameter to RequestContract
		contract := models.RequestContract{
			ID:            i + 1, // Or some unique ID generation logic
			EndpointID:    endpointID,
			ParameterName: param.ParameterName,
			DataType:      param.DataType,
			Description:   param.Description,
			IsRequired:    param.IsRequired,
			DefaultValue:  fmt.Sprintf("%v", param.DefaultValue), // Convert DefaultValue to string
		}

		contracts = append(contracts, contract)
	}

	return contracts
}

func mapParametersToResponseContracts(parameters []Parameter, endpointID int) []models.ResponseContract {
	var contracts []models.ResponseContract

	for i, param := range parameters {
		// Map each Parameter to ResponseContract
		contract := models.ResponseContract{
			ID:            i + 1, // Or some unique ID generation logic
			EndpointID:    endpointID,
			ParameterName: param.ParameterName,
			DataType:      param.DataType,
			Description:   param.Description,
			IsRequired:    param.IsRequired,
			DefaultValue:  fmt.Sprintf("%v", param.DefaultValue), // Convert DefaultValue to string
		}

		contracts = append(contracts, contract)
	}

	return contracts
}

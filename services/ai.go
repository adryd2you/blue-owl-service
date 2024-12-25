package services

import (
	"blue-owl-service/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Define a struct for the form data
type FormData struct {
	ApiType     string                 `json:"apiType"`
	ApiEndpoint string                 `json:"apiEndpoint"`
	HttpMethod  string                 `json:"httpMethod"`
	Request     map[string]interface{} `json:"request"`
	Response    map[string]interface{} `json:"response"`
}

// Define a struct for the message format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Define a struct for the request body
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Define a struct for the API response
type ApiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Parameter struct {
	ParameterName string      `json:"parameter_name"`
	DataType      string      `json:"data_type"`
	Description   string      `json:"description"`
	IsRequired    bool        `json:"is_required"`
	DefaultValue  interface{} `json:"default_value"`
}

// Function to submit API details to the external service and get the documentation
func RequestContractAI(formData models.CreateAPIEndpoint) ([]Parameter, error) {
	// Create the message body for the request
	message := RequestBody{
		Model: "meta-llama/llama-3.2-90b-vision-instruct:free",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Please generate Request Body Description Containing parameter_name, data_type, description, is_required, default_value",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"Endpoint: %s, HTTP Method: %s, Request: %s",
					formData.Path,
					formData.Method,
					formData.RequestBody,
				),
			},
			{
				Role: "user",
				Content: `only in json format with list inside, for example:
			[
        	{
            "parameter_name": "user_id",
            "data_type": "int",
            "description": "Unique identifier for the user",
            "is_required": true,
            "default_value": null
        	},
			{
            "parameter_name": "Room Number",
            "data_type": "int",
            "description": "Information of Room Number",
            "is_required": true,
            "default_value": null
        	}
    		]`,
			},
			{
				Role:    "user",
				Content: "do not add anything else",
			},
		},
	}

	// Marshal the request body into JSON
	reqBody, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Get the API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("API key is missing")
	}

	// Make the POST request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate documentation: status code %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response JSON
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	// Extract the documentation from the response
	documentation := apiResponse.Choices[0].Message.Content
	if documentation == "" {
		return nil, errors.New("documentation content is empty")
	}

	var parameters []Parameter
	err = json.Unmarshal([]byte(documentation), &parameters)
	if err != nil {
		// Handle JSON unmarshalling error
		return nil, errors.New("error unmarshalling documentation")
	}
	return parameters, nil
}

func ResponseContractAI(formData models.CreateAPIEndpoint) ([]Parameter, error) {
	// Create the message body for the request
	if formData.ResponseBody != nil {
		var responseBodyString string
		err := json.Unmarshal([]byte(formData.ResponseBody), &responseBodyString)
		if err != nil {
			log.Printf("error unmarshalling response_body: %v", err)
		}
		formData.ResponseBody = json.RawMessage(responseBodyString)
	}
	cleanedResponse := string(formData.ResponseBody)
	fmt.Print(cleanedResponse)

	cleanedResponse = cleanedResponse[1 : len(cleanedResponse)-1]
	fmt.Print(cleanedResponse)

	cleanedResponse = strings.Replace(cleanedResponse, `\"`, `"`, -1)
	fmt.Print(cleanedResponse)

	cleanedResponse = strings.Replace(cleanedResponse, `\\n`, "", -1)
	fmt.Print(cleanedResponse)
	message := RequestBody{
		Model: "meta-llama/llama-3.2-90b-vision-instruct:free",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Please generate Response Body Description Containing parameter_name, data_type, description, is_required, default_value",
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"Endpoint: %s, HTTP Method: %s, Response: %s",
					formData.Path,
					formData.Method,
					cleanedResponse,
				),
			},
			{
				Role: "user",
				Content: `only in json format with list inside, for example:
			[
        	{
            "parameter_name": "user_id",
            "data_type": "int",
            "description": "Unique identifier for the user",
            "is_required": true,
            "default_value": null
        	},
			{
            "parameter_name": "Room Number",
            "data_type": "int",
            "description": "Information of Room Number",
            "is_required": true,
            "default_value": null
        	}
    		]`,
			},
			{
				Role:    "user",
				Content: "do not add anything else",
			},
		},
	}

	// Marshal the request body into JSON
	reqBody, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Get the API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("API key is missing")
	}

	// Make the POST request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate documentation: status code %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response JSON
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	// Extract the documentation from the response
	documentation := apiResponse.Choices[0].Message.Content
	if documentation == "" {
		return nil, errors.New("documentation content is empty")
	}

	var parameters []Parameter
	err = json.Unmarshal([]byte(documentation), &parameters)
	if err != nil {
		// Handle JSON unmarshalling error
		return nil, errors.New("error unmarshalling documentation")
	}
	return parameters, nil
}

func HitAI(c *gin.Context) {

	// Example usage
	formData := models.CreateAPIEndpoint{
		Path:         "/users",
		Method:       "GET",
		RequestBody:  json.RawMessage(`{"user_id": "123", "phonenumber": "081319333212"}`),
		ResponseBody: json.RawMessage(`{"user_id": "123"}`),
	}

	documentation, err := RequestContractAI(formData)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		c.JSON(http.StatusOK, gin.H{"data": documentation})
	}
}

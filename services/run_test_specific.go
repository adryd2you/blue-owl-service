package services

import (
	"blue-owl-service/models"
	"blue-owl-service/repositories"
	"blue-owl-service/transport"
	"blue-owl-service/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RunSpecificTest(c *gin.Context) {
	var runRequest models.TestRequest

	if runRequest.ProjectID == 0 {
		if err := c.BindJSON(&runRequest); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	example, err := repositories.GetAPIExampleByID(transport.Db, runRequest.ExampleID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	c.IndentedJSON(http.StatusCreated, fmt.Sprintf("Success Run Specific Test by Example ID: %d", example.ID))
}

func ExecuteRunSpecificTest(req *models.TestRequest) error {
	apiExample, err := repositories.GetAPIExampleByID(transport.Db, req.ExampleID)
	if err != nil {
		return err
	}

	parsed, err := AdaptAPIExample(*apiExample)
	if err != nil {
		return err
	}

	go func() {
		testResult, err := ExecuteTest(&parsed)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = repositories.InsertTestResult(transport.Db, testResult)
		if err != nil {
			log.Println(err)
			return
		}

		if testResult.IsSuccess {
			err = repositories.UpdateTestState(transport.Db, req.ExampleID, 1)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			err = repositories.UpdateTestState(transport.Db, req.ExampleID, 2)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	return nil
}

func ExecuteTest(req *models.APIExample) (models.TestResult, error) {
	reqBody, err := json.Marshal(req.RequestBody)
	if err != nil {
		log.Fatalf("Error marshalling APIExample struct: %v", err)
	}
	// Create the HTTP request
	request, err := http.NewRequest(req.APIMethod, req.Endpoint+req.Path, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}

	// Set the Content-Type header to application/json
	request.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	// Print the status code and response
	fmt.Printf("Response status: %s\n", resp.Status)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error ReadAll Body")
		return models.TestResult{}, err
	}
	defer resp.Body.Close()

	// Convert headers to string format
	var headerBuffer bytes.Buffer
	for key, values := range resp.Header {
		headerBuffer.WriteString(key + ": " + values[0] + "\n")
	}

	responseHeader := headerBuffer.String()
	responseBody := string(bodyBytes)
	testResult := models.TestResult{
		ExampleID:          req.ID,
		TestExecutedAt:     time.Now(),
		ResponseStatusCode: resp.StatusCode,
		ResponseHeader:     responseHeader,
		ResponseBody:       responseBody,
	}

	testResult.IsSuccess, testResult.TestMessage = utils.VerifyResponse(*req, testResult)

	return testResult, nil
}

func AdaptAPIExample(unParsed models.APIExampleUnParsed) (models.APIExample, error) {
	// Prepare the target struct
	var parsed models.APIExample
	parsed.ID = unParsed.ID
	parsed.ApiEndpointID = unParsed.ApiEndpointID
	parsed.Endpoint = unParsed.Endpoint
	parsed.Path = unParsed.Path
	parsed.APIMethod = unParsed.APIMethod
	parsed.RequestHeader = unParsed.RequestHeader
	parsed.RequestParameter = unParsed.RequestParameter
	parsed.ResponseStatusCode = unParsed.ResponseStatusCode
	parsed.ResponseHeader = unParsed.ResponseHeader
	parsed.TestState = unParsed.TestState
	parsed.CreatedAt = unParsed.CreatedAt
	parsed.Name = unParsed.Name

	// Unmarshal ResponseBody string into a map
	if unParsed.ResponseBody != "" {
		var responseBodyString string
		var responseMap map[string]interface{}
		err := json.Unmarshal([]byte(unParsed.ResponseBody), &responseBodyString)
		if err != nil {
			return parsed, fmt.Errorf("error unmarshalling response_body: %v", err)
		}

		err = json.Unmarshal([]byte(responseBodyString), &responseMap)
		if err != nil {
			return parsed, fmt.Errorf("error unmarshalling response_body: %v", err)
		}
		parsed.ResponseBody = responseMap
	}

	// Unmarshal ResponseBody string into a map
	if unParsed.RequestBody != "" && unParsed.RequestBody != "{}" && unParsed.RequestBody != "\"\"" {
		var requestBodyString string
		var requestBodyMap map[string]interface{}
		err := json.Unmarshal([]byte(unParsed.RequestBody), &requestBodyString)
		if err != nil {
			return parsed, fmt.Errorf("error unmarshalling request_body: %v", err)
		}
		err = json.Unmarshal([]byte(requestBodyString), &requestBodyMap)
		if err != nil {
			return parsed, fmt.Errorf("error unmarshalling request_body: %v", err)
		}
		parsed.RequestBody = requestBodyMap
	} else {
		parsed.RequestBody = map[string]interface{}{}
	}

	return parsed, nil
}

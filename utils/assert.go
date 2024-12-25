package utils

import (
	"blue-owl-service/models"
	"encoding/json"
	"fmt"
	"reflect"
)

func CheckEqual(expected, actual interface{}) (bool, string) {
	if expected != actual {
		return false, fmt.Sprintf("expected %v, got %v", expected, actual)
	}
	return true, ""
}

func VerifyResponse(expected models.APIExample, actual models.TestResult) (bool, string) {
	// Define the checks array with logic to skip certain checks
	checks := []struct {
		expected    interface{}
		actual      interface{}
		skip        bool
		testMessage string
	}{
		{expected.ResponseStatusCode, actual.ResponseStatusCode, false, "Response Status Code does not match"},
		// {expected.ResponseHeader, actual.ResponseHeader, expected.ResponseHeader == "", "Response Header does not match"},
		{expected.ResponseBody, actual.ResponseBody, len(expected.ResponseBody) == 0, "Response Body does not match"},
	}

	// Iterate through checks and assert equality
	for _, check := range checks {
		// Skip if condition is true
		if check.skip {
			continue
		}

		// Handle ResponseBody comparison if it's a JSON string
		if check.testMessage == "Response Body does not match" {
			// First, ensure that both expected and actual are of the same type

			// If expected is a map, check if actual is a string (JSON) and unmarshal it
			if expectedMap, ok := check.expected.(map[string]interface{}); ok {
				if actualStr, ok := check.actual.(string); ok {
					// If actual is a JSON string, unmarshal it
					var actualParsed map[string]interface{}
					if err := json.Unmarshal([]byte(actualStr), &actualParsed); err != nil {
						return false, fmt.Sprintf("%s: Error unmarshalling actual JSON: %v", check.testMessage, err)
					}
					// Compare the expected map with the parsed actual map
					if !reflect.DeepEqual(expectedMap, actualParsed) {
						return false, fmt.Sprintf("%s: Expected %v, but got %v", check.testMessage, expectedMap, actualParsed)
					}
					continue
				}
			}

			// If expected is a string, unmarshal it into a map and compare it with actual
			if expectedStr, ok := check.expected.(string); ok {
				var expectedParsed map[string]interface{}
				if err := json.Unmarshal([]byte(expectedStr), &expectedParsed); err != nil {
					return false, fmt.Sprintf("%s: Error unmarshalling expected JSON: %v", check.testMessage, err)
				}

				// If actual is a map, compare it with the unmarshalled expected map
				if actualMap, ok := check.actual.(map[string]interface{}); ok {
					if !reflect.DeepEqual(expectedParsed, actualMap) {
						return false, fmt.Sprintf("%s: Expected %v, but got %v", check.testMessage, expectedParsed, actualMap)
					}
					continue
				}
			}
		}

		// For other types, use direct comparison
		if check.expected != check.actual {
			return false, fmt.Sprintf("%s: Expected %v, but got %v", check.testMessage, check.expected, check.actual)
		}
	}

	return true, ""
}

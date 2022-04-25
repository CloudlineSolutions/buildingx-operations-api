package buildingx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Verb string

const (
	GET   Verb = "GET"
	POST  Verb = "POST"
	PATCH Verb = "PATCH"
)

type APIRequest struct {
	Partition string
	JWT       string
	Path      string
	Operation Verb
	Body      bytes.Reader
}
type SBResponse struct {
	Errors []SBErrorResponse `json:"errors"`
}
type SBErrorResponse struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// MakeRESTCall encapsulates a REST API call and returns the results of the call
func MakeRESTCall(apiReq APIRequest) ([]byte, error) {

	result := make([]byte, 0)

	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return result, errors.New("missing buildingx api endpoint")
	}

	url := fmt.Sprintf("%s/operations/partitions/%s/%s", endpoint, apiReq.Partition, apiReq.Path)
	auth := fmt.Sprintf("Bearer %s", apiReq.JWT)
	client := &http.Client{Timeout: time.Duration(20) * time.Second}
	req, _ := http.NewRequest(string(apiReq.Operation), url, &apiReq.Body)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", auth)

	if apiReq.Operation == PATCH {
		req.Header.Add("content-type", "application/vnd.api+json")
	} else {
		req.Header.Add("content-type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("unexpected error while invoking http client: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		// attempt to parse the error response message
		errorResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, fmt.Errorf("unexpected error trying to read API error response message: %s", err.Error())
		}
		sbResponse := SBResponse{}
		if err := json.Unmarshal(errorResponse, &sbResponse); err != nil {
			return result, fmt.Errorf("Error parsing API response: %s", err.Error())
		}
		if len(sbResponse.Errors) < 1 {
			return result, fmt.Errorf("the building x API returned a status code of %s but no error message was included", resp.Status)
		}

		// the error response message was successfully parsed, so use it to construct an error message
		return result, fmt.Errorf("the building x API returned a status code of %s and an error detail message as follows: %s", sbResponse.Errors[0].Status, sbResponse.Errors[0].Detail)

	}

	//all is well, return the response payload
	return ioutil.ReadAll(resp.Body)

}

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
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", auth)

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			return result, errors.New("error message from BuildingX: " + data["detail"].(string))
		}
		return result, errors.New("got non-200 response from BuildingX API with no additional information")

	}

	return ioutil.ReadAll(resp.Body)

}

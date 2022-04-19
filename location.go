package buildingx

import (
	"encoding/json"
	"errors"
	"os"
)

type Location struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type SBLocationResponse struct {
	Locations []SBLocation `json:"data"`
}
type SBLocation struct {
	ID         string               `json:"id"`
	Attributes SBLocationAttributes `json:"attributes"`
}
type SBLocationAttributes struct {
	TimeZone    string `json:"timeZone"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

func GetLocations(jwt, partitionID string) ([]Location, error) {

	locations := make([]Location, 0)

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return locations, errors.New("missing buildingx api endpoint")
	}

	req := APIRequest{
		Partition: partitionID,
		JWT:       jwt,
		Path:      "locations?filter[type]=Building",
		Verb:      "GET",
	}
	resp, err := MakeRESTCall(req)
	if err != nil {
		return locations, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native location response payload
	sbLocationResponse := SBLocationResponse{}

	if err := json.Unmarshal(resp, &sbLocationResponse); err != nil {
		return locations, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	// now create the Location objects
	for _, sbLocation := range sbLocationResponse.Locations {

		location := Location{
			ID:   sbLocation.ID,
			Name: sbLocation.Attributes.Label,
		}
		locations = append(locations, location)

	}

	return locations, nil

}

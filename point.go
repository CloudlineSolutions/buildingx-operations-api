package buildingx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Point struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DataType    string `json:"dataType"`
	Status      string `json:"status"`
	StringValue string `json:"stringValue"`
	TimeStamp   string `json:"timestamp"`
}
type SBPointResponse struct {
	Points []SBPoint `json:"data"`
}
type SBPoint struct {
	ID         string            `json:"id"`
	Attributes SBPointAttributes `json:"attributes"`
}
type SBPointAttributes struct {
	Name             string                  `json:"name"`
	DataType         string                  `json:"dataType"`
	SystemAttributes SBPointSystemAttributes `json:"systemAttributes"`
	PointValue       SBPointValue            `json:"pointValue"`
}
type SBPointSystemAttributes struct {
	CurStatus   string `json:"curStatus"`
	Description string `json:"description"`
}
type SBPointValue struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

// returns an array of points that are associated with a particular device
func GetPointsByDevice(session *Session, deviceID string) ([]Point, error) {

	points := make([]Point, 0)

	// make sure session is initialized
	if !session.IsInitialized {
		return points, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return points, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("devices/%s/points?field[Point]=pointValue", deviceID)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return points, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native points response payload
	sbPointsResponse := SBPointResponse{}
	if err := json.Unmarshal(resp, &sbPointsResponse); err != nil {
		return points, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	for _, sbPoint := range sbPointsResponse.Points {
		point := Point{
			ID:          sbPoint.ID,
			Name:        sbPoint.Attributes.Name,
			Description: sbPoint.Attributes.SystemAttributes.Description,
			DataType:    sbPoint.Attributes.DataType,
			Status:      sbPoint.Attributes.SystemAttributes.CurStatus,
			StringValue: sbPoint.Attributes.PointValue.Value,
			TimeStamp:   sbPoint.Attributes.PointValue.Timestamp,
		}

		points = append(points, point)
	}

	return points, nil
}

package buildingx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Point struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DataType    string    `json:"dataType"`
	Writable    bool      `json:"writable"`
	Status      string    `json:"status"`
	StringValue string    `json:"stringValue"`
	Timestamp   time.Time `json:"timestamp"`
}
type SBPointsResponse struct {
	Points []SBPoint `json:"data"`
}
type SBPointResponse struct {
	Point SBPoint `json:"data"`
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
	Writable    string `json:"writable"`
}
type SBPointValue struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type PointHistory struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}
type SBPointHistoryResponse struct {
	Data []SBPointHistory `json:"data"`
}
type SBPointHistory struct {
	Attributes SBPointHistoryAttributes `json:"attributes"`
}
type SBPointHistoryAttributes struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type PointHistoryRequest struct {
	Point Point
	Start time.Time
	End   time.Time
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
	sbPointsResponse := SBPointsResponse{}
	if err := json.Unmarshal(resp, &sbPointsResponse); err != nil {
		return points, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	//TODO: create a common point mapping function for this function and GetSinglePoint
	for _, sbPoint := range sbPointsResponse.Points {

		// deliberately ignoring the error here as we don't know what to do with it
		timeStamp, _ := time.Parse(time.RFC3339, sbPoint.Attributes.PointValue.Timestamp)
		writableString := sbPoint.Attributes.SystemAttributes.Writable
		writable := false
		if writableString == "m:" {
			writable = true
		}

		point := Point{
			ID:          sbPoint.ID,
			Name:        sbPoint.Attributes.Name,
			Description: sbPoint.Attributes.SystemAttributes.Description,
			DataType:    sbPoint.Attributes.DataType,
			Writable:    writable,
			Status:      sbPoint.Attributes.SystemAttributes.CurStatus,
			StringValue: sbPoint.Attributes.PointValue.Value,
			Timestamp:   timeStamp,
		}

		points = append(points, point)
	}

	return points, nil
}
func GetSinglePoint(session *Session, id string) (Point, error) {

	point := Point{}
	// make sure session is initialized
	if !session.IsInitialized {
		return point, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return point, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("points/%s?field[Point]=pointValue", id)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return point, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native point response payload
	sbPointResponse := SBPointResponse{}
	if err := json.Unmarshal(resp, &sbPointResponse); err != nil {
		return point, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	// map the native point structure to our point structure

	// deliberately ignoring the error here as we don't know what to do with it
	timeStamp, _ := time.Parse(time.RFC3339, sbPointResponse.Point.Attributes.PointValue.Timestamp)
	writableString := sbPointResponse.Point.Attributes.SystemAttributes.Writable
	writable := false
	if writableString == "m:" {
		writable = true
	}

	point.ID = sbPointResponse.Point.ID
	point.Name = sbPointResponse.Point.Attributes.Name
	point.Description = sbPointResponse.Point.Attributes.SystemAttributes.Description
	point.DataType = sbPointResponse.Point.Attributes.DataType
	point.Writable = writable
	point.Status = sbPointResponse.Point.Attributes.SystemAttributes.CurStatus
	point.StringValue = sbPointResponse.Point.Attributes.PointValue.Value
	point.Timestamp = timeStamp

	// all is well. return the point
	return point, nil

}
func GetPointHistory(session Session, point *Point, start, end time.Time) ([]PointHistory, error) {

	history := make([]PointHistory, 0)

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return history, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("points/%s/values?filter[timestamp][from]=%s&[timestamp][to]=%s", point.ID, start.Format(time.RFC3339), end.Format(time.RFC3339))
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return history, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native point response payload
	sbPointHistoryResponse := SBPointHistoryResponse{}
	if err := json.Unmarshal(resp, &sbPointHistoryResponse); err != nil {
		return history, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	for _, sbHistory := range sbPointHistoryResponse.Data {

		pointHistory := PointHistory{
			Value:     sbHistory.Attributes.Value,
			Timestamp: sbHistory.Attributes.Timestamp,
		}

		history = append(history, pointHistory)

	}

	return history, nil

}

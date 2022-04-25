package buildingx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Location struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Street      string `json:"street"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
	Country     string `json:"country"`
	TimeZone    string `json:"timeZone"`
}
type SBLocationsResponse struct {
	Locations []SBLocation `json:"data"`
}
type SBLocationResponse struct {
	Location SBLocation `json:"data"`
}
type SBLocation struct {
	ID            string                  `json:"id"`
	Attributes    SBLocationAttributes    `json:"attributes"`
	Relationships SBLocationRelationships `json:"relationships"`
}
type SBLocationAttributes struct {
	TimeZone    string `json:"timeZone"`
	Label       string `json:"label"`
	Description string `json:"description"`
}
type SBLocationRelationships struct {
	Features SBLocationHasAddress `json:"hasPostalAddress"`
}
type SBLocationHasAddress struct {
	Data SBLocationPostalAddressData `json:"data"`
}
type SBLocationPostalAddressData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type SBLocationIncludedResponse struct {
	Included []SBLocationIncluded `json:"included"`
}
type SBLocationIncluded struct {
	ID         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes SBLocationIncludedAttributes `json:"attributes"`
}
type SBLocationIncludedAttributes struct {
	Locality      string `json:"locality"`
	CountryCode   string `json:"countryCode"`
	CountryName   string `json:"countryName"`
	ContinentCode string `json:"continentCode"`
	ContinentName string `json:"continentName"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	Street        string `json:"street"`
}

// GetLocations returns an array of all locations associated with the session. It also populates the session object with the locations.
func GetLocations(session *Session) ([]Location, error) {

	locations := make([]Location, 0)

	// make sure session is initialized
	if !session.IsInitialized {
		return locations, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return locations, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      "locations?filter[type]=Building&include=hasPostalAddress",
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return locations, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native location response payload
	sbLocationsResponse := SBLocationsResponse{}
	if err := json.Unmarshal(resp, &sbLocationsResponse); err != nil {
		return locations, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	sbLocationsIncludedResponse := SBLocationIncludedResponse{}
	if err := json.Unmarshal(resp, &sbLocationsIncludedResponse); err != nil {
		return locations, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	// now create the Location objects
	for _, sbLocation := range sbLocationsResponse.Locations {

		location := Location{
			ID:          sbLocation.ID,
			Name:        sbLocation.Attributes.Label,
			Description: sbLocation.Attributes.Description,
			TimeZone:    sbLocation.Attributes.TimeZone,
		}

		for _, include := range sbLocationsIncludedResponse.Included {
			if sbLocation.Relationships.Features.Data.ID == include.ID {
				location.City = include.Attributes.Locality
				location.Street = include.Attributes.Street
				location.Country = include.Attributes.CountryCode
				location.PostalCode = include.Attributes.PostalCode
			}
		}

		locations = append(locations, location)

	}

	// all is well. return the locations
	return locations, nil

}
func GetSingleLocation(session *Session, id string) (Location, error) {

	location := Location{}
	// make sure session is initialized
	if !session.IsInitialized {
		return location, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return location, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("locations/%s?include=hasPostalAddress", id)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return location, errors.New("error making REST call: " + err.Error())
	}

	// Unmarshal the native location response payload
	sbLocationResponse := SBLocationResponse{}
	if err := json.Unmarshal(resp, &sbLocationResponse); err != nil {
		return location, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	sbLocationsIncludedResponse := SBLocationIncludedResponse{}
	if err := json.Unmarshal(resp, &sbLocationsIncludedResponse); err != nil {
		return location, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	location.ID = sbLocationResponse.Location.ID
	location.Description = sbLocationResponse.Location.Attributes.Description
	location.Name = sbLocationResponse.Location.Attributes.Label
	location.TimeZone = sbLocationResponse.Location.Attributes.TimeZone

	for _, include := range sbLocationsIncludedResponse.Included {
		if sbLocationResponse.Location.Relationships.Features.Data.ID == include.ID {
			location.City = include.Attributes.Locality
			location.Street = include.Attributes.Street
			location.Country = include.Attributes.CountryCode
			location.PostalCode = include.Attributes.PostalCode
		}
	}

	// all is well. return the location
	return location, nil

}

package buildingx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Device struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Model        string `json:"model"`
	Serial       string `json:"serial"`
	OnlineStatus string `json:"onlineStatus"`
}
type SBDevicesResponse struct {
	Devices []SBDevice `json:"data"`
}
type SBDeviceResponse struct {
	Device SBDevice `json:"data"`
}
type SBDevice struct {
	ID            string                `json:"id"`
	Attributes    SBDeviceAttributes    `json:"attributes"`
	RelationShips SBDeviceRelationships `json:"relationships"`
}
type SBDeviceRelationships struct {
	Features SBDeviceFeatures `json:"hasFeatures"`
}
type SBDeviceFeatures struct {
	Data []SBDeviceFeaturesData `json:"data"`
}
type SBDeviceFeaturesData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type SBDeviceAttributes struct {
	ModelName    string `json:"modelName"`
	SerialNumber string `json:"serialNumber"`
}
type SBDevicesIncludedResponse struct {
	Included []SBDeviceIncluded `json:"included"`
}
type SBDeviceIncluded struct {
	ID            string                        `json:"id"`
	Type          string                        `json:"type"`
	Attributes    SBDeviceIncludedAttributes    `json:"attributes"`
	RelationShips SBDeviceIncludedRelationships `json:"relationships"`
}
type SBDeviceIncludedAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type SBDeviceIncludedRelationships struct {
	HasDevice SBDeviceIncludedRelationshipsHasDevice `json:"hasDevice"`
}
type SBDeviceIncludedRelationshipsHasDevice struct {
	Data SBDeviceIncludedRelationshipsData `json:"data"`
}
type SBDeviceIncludedRelationshipsData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// returns an array of devices that are associated with a particular location
func GetDevicesByLocation(session *Session, location *Location) ([]Device, error) {

	devices := make([]Device, 0)

	// make sure session is initialized
	if !session.IsInitialized {
		return devices, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return devices, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("devices?include=hasFeatures.DeviceInfo,hasFeatures.Connectivity&filter[hasLocation.data.id]=%s", location.ID)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return devices, errors.New("error making REST call: " + err.Error())
	}

	return parseDevicesJSON(resp)

}

// returns an array of devices that are associated with a particular gateway
func GetDevicesByGateway(session *Session, gatewayID string) ([]Device, error) {

	devices := make([]Device, 0)

	// make sure session is initialized
	if !session.IsInitialized {
		return devices, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return devices, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("devices/%s/devices?include=hasFeatures.DeviceInfo", gatewayID)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return devices, errors.New("error making REST call: " + err.Error())
	}

	return parseDevicesJSON(resp)

}

// returns an array of devices that are associated with the partition
func GetAllDevices(session *Session) ([]Device, error) {

	devices := make([]Device, 0)

	// make sure session is initialized
	if !session.IsInitialized {
		return devices, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return devices, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("devices?include=hasFeatures.DeviceInfo,hasFeatures.Connectivity")
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return devices, errors.New("error making REST call: " + err.Error())
	}

	return parseDevicesJSON(resp)

}
func parseDevicesJSON(payload []byte) ([]Device, error) {

	devices := make([]Device, 0)

	// Unmarshal the native devices response payload
	sbDevicesResponse := SBDevicesResponse{}
	if err := json.Unmarshal(payload, &sbDevicesResponse); err != nil {
		return devices, errors.New("Error parsing API response. String submitted: " + string(payload))
	}

	// Now unmarshal the device features nodes
	sbDevicesIncludedResponse := SBDevicesIncludedResponse{}
	if err := json.Unmarshal(payload, &sbDevicesIncludedResponse); err != nil {
		return devices, errors.New("Error parsing API response (features section). String submitted: " + string(payload))
	}

	// now create the Location objects
	for _, sbDevice := range sbDevicesResponse.Devices {

		device := Device{
			ID:     sbDevice.ID,
			Model:  sbDevice.Attributes.ModelName,
			Serial: sbDevice.Attributes.SerialNumber,
		}
		// loop through the device features to populate the rest of the properties on the Device
		for _, sbFeature := range sbDevicesIncludedResponse.Included {
			// there are two kinds of devices features: DeviceInfo and Connectivity. Get the desired properties from each.
			if sbFeature.RelationShips.HasDevice.Data.ID == device.ID {
				if strings.ToLower(sbFeature.Type) == "deviceinfo" {
					device.Name = sbFeature.Attributes.Name
					device.Description = sbFeature.Attributes.Description
				} else if strings.ToLower(sbFeature.Type) == "connectivity" {
					device.OnlineStatus = sbFeature.Attributes.Status
				}

			}

		}
		if device.OnlineStatus == "" {
			device.OnlineStatus = "Unknown"
		}
		devices = append(devices, device)

	}

	// all is well. return the devices
	return devices, nil

}

// returns a single device by its id
func GetSingleDevice(session *Session, id string) (Device, error) {

	device := Device{}
	// make sure session is initialized
	if !session.IsInitialized {
		return device, errors.New("session is not initialized")
	}

	// make sure you have the required environment variable
	endpoint := os.Getenv("BUILDINGX_ENDPOINT")
	if endpoint == "" {
		return device, errors.New("missing buildingx api endpoint")
	}

	// create the API request
	path := fmt.Sprintf("devices/%s", id)
	req := APIRequest{
		Partition: session.Partition,
		JWT:       session.JWT,
		Path:      path,
		Operation: GET,
	}

	// make the API call
	resp, err := MakeRESTCall(req)
	if err != nil {
		return device, errors.New("error making REST call: " + err.Error())
	}

	//TODO This function has code duplication for parsing device information. consolidate with GetAllDevices, GetDevicesByLocation and GetDevicesByGateway

	// Unmarshal the native device response payload
	sbDeviceResponse := SBDeviceResponse{}
	if err := json.Unmarshal(resp, &sbDeviceResponse); err != nil {
		return device, errors.New("Error parsing API response. String submitted: " + string(resp))
	}

	// Now unmarshal the device features nodes
	sbDevicesIncludedResponse := SBDevicesIncludedResponse{}
	if err := json.Unmarshal(resp, &sbDevicesIncludedResponse); err != nil {
		return device, errors.New("Error parsing API response (features section). String submitted: " + string(resp))
	}

	device.ID = sbDeviceResponse.Device.ID
	device.Model = sbDeviceResponse.Device.Attributes.ModelName
	device.Serial = sbDeviceResponse.Device.Attributes.SerialNumber

	// loop through the device features to populate the rest of the properties on the Device
	for _, sbFeature := range sbDevicesIncludedResponse.Included {
		// there are two kinds of devices features: DeviceInfo and Connectivity. Get the desired properties from each.
		if sbFeature.RelationShips.HasDevice.Data.ID == device.ID {
			if strings.ToLower(sbFeature.Type) == "deviceinfo" {
				device.Name = sbFeature.Attributes.Name
				device.Description = sbFeature.Attributes.Description
			} else if strings.ToLower(sbFeature.Type) == "connectivity" {
				device.OnlineStatus = sbFeature.Attributes.Status
			}

		}

	}
	if device.OnlineStatus == "" {
		device.OnlineStatus = "Unknown"
	}

	// all is well. return the device
	return device, nil

}

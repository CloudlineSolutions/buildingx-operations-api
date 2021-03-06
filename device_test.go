package buildingx

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/stretchr/testify/assert"
)

func TestGetDevicesByLocation(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetDevicesByLocation")

	// first make sure you have a partition ID (from environment variable)
	partitionID := os.Getenv("BUILDINGX_PARTITION_ID")
	if partitionID == "" {
		t.Fatal("unable to find partition ID in environment variable")
	}

	// initialize the session (uses credentials to authenticate and produce a JWT)
	session := Session{}
	err := session.Initialize(partitionID)
	if err != nil {
		t.Fatal("test failed while initializing session: ", err.Error())
	}

	// get the location collection to use in fetching a single device by location
	// If there are no locations, this will fail
	locations, err := GetLocations(&session)
	if err != nil {
		t.Fatal("error getting locations: ", err.Error())
	}
	if len(locations) < 1 {
		t.Fatal("test failed because there are no valid locations to use")
	}

	t.Run("get-devices-with-valid-location-id", func(t *testing.T) {
		devices, err := GetDevicesByLocation(&session, &locations[0])
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// we should have gotten at least one device
		assert.GreaterOrEqual(t, len(devices), 1)
	})
	t.Run("get-devices-with-invalid-location-id", func(t *testing.T) {

		badLocation := Location{
			ID: "invalid-id",
		}
		devices, err := GetDevicesByLocation(&session, &badLocation)
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// we should have gotten no devices
		assert.Equal(t, 0, len(devices))
	})

}
func TestGetDevicesByGateway(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetDevicesByGateway")

	// first make sure you have a partition ID (from environment variable)
	partitionID := os.Getenv("BUILDINGX_PARTITION_ID")
	if partitionID == "" {
		t.Fatal("unable to find partition ID in environment variable")
	}

	// initialize the session (uses credentials to authenticate and produce a JWT)
	session := Session{}
	err := session.Initialize(partitionID)
	if err != nil {
		t.Fatal("test failed while initializing session: ", err.Error())
	}

	// Get all Devices associated with the partion and find the gateway.
	// If a gateway (X300 or X200) is not present, this will fail
	devices, err := GetAllDevices(&session)
	if err != nil {
		t.Fatal("error getting all devices: ", err.Error())
	}
	if len(devices) < 1 {
		t.Fatal("test failed because there are no valid devices to use")
	}

	// find the gateway
	gatewayID := ""
	for _, device := range devices {
		if strings.ToLower(device.Model) == "x300" || strings.ToLower(device.Model) == "x200" {
			gatewayID = device.ID
			break
		}
	}
	if gatewayID == "" {
		t.Fatal("could not find a valid gateway in the device collection")
	}

	t.Run("get-devices-with-valid-gateway-id", func(t *testing.T) {
		devices, err := GetDevicesByGateway(&session, gatewayID)
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// we should have gotten at least one device
		assert.GreaterOrEqual(t, len(devices), 1)
	})
	t.Run("get-devices-with-invalid-gateway-id", func(t *testing.T) {
		gatewayID = "invalid-id"
		_, err := GetDevicesByGateway(&session, gatewayID)

		// The Building X API returns a 404 error when the device ID is invalid
		assert.NotNil(t, err)
	})

}
func TestGetSingleDevice(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetSingleDevice")

	// first make sure you have a partition ID (from environment variable)
	partitionID := os.Getenv("BUILDINGX_PARTITION_ID")
	if partitionID == "" {
		t.Fatal("unable to find partition ID in environment variable")
	}

	// initialize the session (uses credentials to authenticate and produce a JWT)
	session := Session{}
	err := session.Initialize(partitionID)
	if err != nil {
		t.Fatal("test failed while initializing session: ", err.Error())
	}

	// Get all Devices associated with the partion and use the first one as a reference
	devices, err := GetAllDevices(&session)
	if err != nil {
		t.Fatal("error getting all devices: ", err.Error())
	}
	if len(devices) < 1 {
		t.Fatal("test failed because there are no valid devices to use")
	}

	t.Run("get-single-device", func(t *testing.T) {
		device, err := GetSingleDevice(&session, devices[0].ID)
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// an empty ID on the Location object means that one was not found
		assert.Equal(t, devices[0].ID, device.ID)
	})
	t.Run("get-single-device-with-invalid-id", func(t *testing.T) {
		deviceID := "invalid-id"
		_, err := GetSingleDevice(&session, deviceID)

		// The Building X API returns a 404 error when the device ID is invalid
		assert.NotNil(t, err)
	})
}

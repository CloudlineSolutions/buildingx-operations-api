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

	// get the location collection to use in fetching a single valid location
	// If there are no locations, this will fail
	locations, err := GetLocations(&session)
	if err != nil {
		t.Fatal("error getting locations: ", err.Error())
	}
	if len(locations) < 1 {
		t.Fatal("test failed because there are no valid locations to use")
	}

	t.Run("get-devices-with-valid-location-id", func(t *testing.T) {
		devices, err := GetDevicesByLocation(&session, locations[0].ID)
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// an empty ID on the Location object means that one was not found
		assert.GreaterOrEqual(t, 1, len(devices))
	})
	//TODO add test with invalid location
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

	// get all of the devices associated with a location
	t.Run("get-devices-with-valid-gateway-id", func(t *testing.T) {
		devices, err := GetDevicesByGateway(&session, gatewayID)
		if err != nil {
			t.Fatal("error getting devices: ", err.Error())
		}
		// an empty ID on the Location object means that one was not found
		assert.GreaterOrEqual(t, 1, len(devices))
	})
	//TODO add test with invalid gateway
}

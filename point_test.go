package buildingx

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/stretchr/testify/assert"
)

func TestGetPoints(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetPoints")

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

	// get the first device under a gateway. we expect points on this one.
	gatewayDevices, err := GetDevicesByGateway(&session, gatewayID)
	if err != nil {
		t.Fatal("error getting devices under the gateway: ", err.Error())
	}

	if len(gatewayDevices) < 1 {
		t.Fatal("no devices found under the gateway")
	}

	// get all of the points associated with the device under a gateway
	t.Run("get-points-with-valid-device-id", func(t *testing.T) {
		points, err := GetPointsByDevice(&session, gatewayDevices[0].ID)
		if err != nil {
			t.Fatal("error getting points: ", err.Error())
		}
		// an empty ID on the Location object means that one was not found
		assert.GreaterOrEqual(t, 1, len(points))
	})
	//TODO add test with invalid device id
}

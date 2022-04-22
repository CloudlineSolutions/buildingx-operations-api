package buildingx

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

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
		assert.GreaterOrEqual(t, len(points), 1)
	})
	//TODO add test with invalid device id
}
func TestGetSinglePoint(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetSinglePoint")

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

	// now get the points
	points, err := GetPointsByDevice(&session, gatewayDevices[0].ID)
	if err != nil {
		t.Fatal("error getting points: ", err.Error())
	}

	if len(points) < 1 {
		t.Fatal("no points found to use in test")
	}

	// get all of the points associated with the device under a gateway
	t.Run("get-points-with-valid-device-id", func(t *testing.T) {
		point, err := GetSinglePoint(&session, points[0].ID)
		if err != nil {
			t.Fatal("error getting point: ", err.Error())
		}
		// make sure the id of the point retrieved is the same as what was asked for
		assert.Equal(t, points[0].ID, point.ID)
	})
	//TODO add test with invalid device id
}
func TestGetPointHistory(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetPointHistory")

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

	// now get the points
	points, err := GetPointsByDevice(&session, gatewayDevices[0].ID)
	if err != nil {
		t.Fatal("error getting points: ", err.Error())
	}

	if len(points) < 1 {
		t.Fatal("no points found to use in test")
	}

	// get all of the points associated with the device under a gateway
	t.Run("get-point-history", func(t *testing.T) {

		// retrieve from 30 days ago
		start := time.Now().UTC().Add(-720 * time.Hour)

		history, err := GetPointHistory(session, &points[0], start, time.Now().UTC())
		if err != nil {
			t.Fatal("error getting point: ", err.Error())
		}
		// make sure the id of the point retrieved is the same as what was asked for
		assert.GreaterOrEqual(t, len(history), 1)
	})
	//TODO add test with invalid device id
}

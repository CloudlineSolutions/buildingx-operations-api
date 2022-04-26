package buildingx

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/stretchr/testify/assert"
)

func TestGetLocations(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetLocations")

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

	t.Run("get-locations-with-valid-partition", func(t *testing.T) {
		locations, err := GetLocations(&session)
		if err != nil {
			t.Fatal("error getting locations: ", err.Error())
		}
		// we should get at least one location
		assert.GreaterOrEqual(t, 1, len(locations))
	})
	t.Run("get-locations-with-invalid-partition", func(t *testing.T) {
		s := Session{
			IsInitialized: true,
			JWT:           session.JWT,
			Partition:     "invalid-partition",
		}
		// we just test here to make sure there is an error
		_, err := GetLocations(&s)
		assert.NotNil(t, err)
	})

}
func TestGetSingleLocation(t *testing.T) {
	ctx := context.Background()
	ctx, _ = xray.BeginSegment(ctx, "TestGetSingleLocations")

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

	t.Run("get-single-location-with-valid-location-id", func(t *testing.T) {
		location, err := GetSingleLocation(&session, locations[0].ID)
		if err != nil {
			t.Fatal("error getting locations: ", err.Error())
		}
		// an empty ID on the Location object means that one was not found
		assert.NotEmpty(t, location.ID)
	})
	t.Run("get-locations-with-invalid-location-id", func(t *testing.T) {

		// we just test here to make sure there is an error
		_, err := GetSingleLocation(&session, "invalid-location-id")
		assert.NotNil(t, err)
	})

}

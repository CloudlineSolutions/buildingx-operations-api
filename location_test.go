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
	jwt, err := GetToken()
	if err != nil {
		t.Fatal("test failed while getting jwt: ", err.Error())
	}

	t.Run("get-locations", func(t *testing.T) {

		locations, err := GetLocations(jwt, partitionID)
		if err != nil {
			t.Fatal("error while retrieving locations: ", err.Error())
		}
		assert.GreaterOrEqual(t, 1, len(locations))
	})

}

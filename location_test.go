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

	// locations are loaded during session initialization so thats all we need to do
	session := Session{}
	err := session.Initialize(partitionID)
	if err != nil {
		t.Fatal("test failed while initializing session: ", err.Error())
	}

	t.Run("get-locations-with-valid-partition", func(t *testing.T) {

		// we expect that there is at least one location
		assert.GreaterOrEqual(t, 1, len(session.Locations))
	})

	//TODO: add test for invalid partition

}

package tracking51

import (
	"fmt"
	"testing"
)

func TestTrackingService_All(t *testing.T) {
	items, _, err := client.Services.Tracking.All(TracksQueryParams{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(items)
}

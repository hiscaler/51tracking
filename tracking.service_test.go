package tracking51

import (
	"fmt"
	"testing"
)

func TestTrackingService_Query(t *testing.T) {
	items, _, err := client.Services.Tracking.Query(TracksQueryParams{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(items)
}

func TestTrackingService_Refresh(t *testing.T) {
	req := RefreshTrackRequests{
		RefreshTrackRequest{
			TrackingNumber: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			CourierCode:    "cms",
		},
	}
	successItems, errorItems, err := client.Services.Tracking.Refresh(req)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v, %#v", successItems, errorItems)
}

func TestTrackingService_StatusStatistic(t *testing.T) {
	stat, err := client.Services.Tracking.StatusStatistic(StatusStatisticRequest{})
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", stat)
}

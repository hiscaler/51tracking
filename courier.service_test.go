package tracking51

import (
	"testing"
)

func TestCarrierService_List(t *testing.T) {
	_, err := client.Services.Courier.List("cn")
	if err != nil {
		t.Error(err)
	}
}

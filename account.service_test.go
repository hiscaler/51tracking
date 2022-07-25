package tracking51

import "testing"

func TestAccountService_Profile(t *testing.T) {
	profile, err := client.Services.Account.Profile()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", profile)
}

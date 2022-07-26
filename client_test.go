package tracking51

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/51tracking-go/config"
	"os"
	"testing"
)

var client *Tracking51

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	client = NewTracking51(c)
	client.SetDebug(c.Debug)
	m.Run()
}

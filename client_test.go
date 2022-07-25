package tracking51

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var client *Tracking51

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	c := struct {
		Debug   bool
		Version string
		AppKey  string
	}{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	client = NewTracking51(c.AppKey)
	client.SetDebug(c.Debug)
	m.Run()
}

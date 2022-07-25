package tracking51

import (
	"encoding/json"
	"github.com/hiscaler/gox/inx"
	"gopkg.in/guregu/null.v4"
)

type carrierService service

type Carrier1 struct {
	Name        string      `json:"courier_name"`
	Code        string      `json:"courier_code"`
	Phone       string      `json:"courier_phone"`
	CountryCode null.String `json:"country_code"`
	Type        string      `json:"courier_type"`
	URL         null.String `json:"courier_url"`
	Logo        string      `json:"courier_logo"`
}
type Carrier struct {
	Name        string      `json:"Name"`
	Express     string      `json:"express"`
	Phone       string      `json:"phone"`
	CountryCode null.String `json:"country_code"`
	TrackURL    null.String `json:"track_url"`
	Picture     string      `json:"picture"`
}

func (s carrierService) List(lang string) (items []Carrier, err error) {
	if !inx.StringIn(lang, "cn", "en") {
		lang = "cn"
	}

	res := struct {
		NormalResponse
		Data []Carrier `json:"data"`
	}{}
	resp, err := s.httpClient.R().
		SetQueryParam("lang", lang).
		Get("/trackings/carriers")
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		items = res.Data
	}
	return
}

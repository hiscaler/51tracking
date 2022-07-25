package tracking51

import (
	"encoding/json"
	"github.com/hiscaler/gox/inx"
	"gopkg.in/guregu/null.v4"
)

type carrierService service

// https://www.51tracking.com/v3/api-index?language=Golang#%E7%89%A9%E6%B5%81%E5%95%86%E5%88%97%E8%A1%A8
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

// List 物流商列表
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

// Change 修改物流简码
func (s carrierService) Change(trackingNumber, oldCourierCode, newCourierCode string) error {
	_, err := s.httpClient.R().
		SetBody(map[string]string{
			"tracking_number":  trackingNumber,
			"courier_code":     oldCourierCode,
			"new_courier_code": newCourierCode,
		}).
		Put("/trackings/modifycourier")
	return err
}

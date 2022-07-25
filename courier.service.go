package tracking51

import (
	"encoding/json"
	"github.com/hiscaler/gox/inx"
	"gopkg.in/guregu/null.v4"
)

type courierService service

// https://www.51tracking.com/v3/api-index?language=Golang#%E7%89%A9%E6%B5%81%E5%95%86%E5%88%97%E8%A1%A8
type Courier struct {
	Name        string      `json:"courier_name"`
	Code        string      `json:"courier_code"`
	Phone       string      `json:"courier_phone"`
	CountryCode null.String `json:"country_code"`
	Type        string      `json:"courier_type"`
	URL         null.String `json:"courier_url"`
	Logo        string      `json:"courier_logo"`
}

// List 物流商列表
func (s courierService) List(lang string) (items []Courier, err error) {
	if !inx.StringIn(lang, "cn", "en") {
		lang = "cn"
	}

	res := struct {
		NormalResponse
		Data []Courier `json:"data"`
	}{}
	resp, err := s.httpClient.R().
		SetQueryParam("lang", lang).
		Get("/trackings/courier")
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		items = res.Data
	}
	return
}

// Change 修改物流简码
func (s courierService) Change(trackingNumber, oldCourierCode, newCourierCode string) error {
	_, err := s.httpClient.R().
		SetBody(map[string]string{
			"tracking_number":  trackingNumber,
			"courier_code":     oldCourierCode,
			"new_courier_code": newCourierCode,
		}).
		Put("/trackings/modifycourier")
	return err
}

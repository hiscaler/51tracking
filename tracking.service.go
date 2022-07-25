package tracking51

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

type trackingService service

type CreateTrackingRequest struct {
	TrackingNumber          string `json:"tracking_number"`                     //	包裹物流单号
	CourierCode             string `json:"courier_code"`                        //	物流商对应的唯一简码
	OrderNumber             string `json:"order_number,omitempty"`              // 包裹的订单号，由商家/平台所产生的订单编号
	Title                   string `json:"title,omitempty"`                     // 包裹名称
	DestinationCode         string `json:"destination_code,omitempty"`          // 目的国的二字简码
	LogisticsChannel        string `json:"logistics_channel,omitempty"`         // 自定义字段，用于填写物流渠道（比如某货代商）
	Note                    string `json:"note,omitempty"`                      // 备注
	CustomerName            string `json:"customer_name,omitempty"`             // 客户姓名
	CustomerEmail           string `json:"customer_email,omitempty"`            // 客户邮箱
	CustomerPhone           string `json:"customer_phone,omitempty"`            // 顾客接收短信的手机号码。手机号码的格式应该为：“+区号手机号码”（例子：+8612345678910）
	ShippingDate            string `json:"shipping_date,omitempty"`             // 包裹发货时间（例子：2020-09-17 16:51）
	TrackingShippingDate    string `json:"tracking_shipping_date,omitempty"`    // 包裹的发货时间，其格式为：YYYYMMDD，有部分的物流商（如 deutsch-post）需要这个参数（例子：20200102）
	TrackingPostalCode      string `json:"tracking_postal_code,omitempty"`      // 收件人所在地邮编，仅有部分的物流商（如 postnl-3s）需要这个参数
	TrackingDestinationCode string `json:"tracking_destination_code,omitempty"` // 目的国对应的二字简码，部分物流商（如postnl-3s）需要这个参数
	TrackingCourierAccount  string `json:"tracking_courier_account,omitempty"`  // 物流商的官方账号，仅有部分的物流商（如 dynamic-logistics）需要这个参数
}

func (m CreateTrackingRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
		validation.Field(&m.CourierCode, validation.Required.Error("物流商简码不能为空")),
		validation.Field(&m.CustomerEmail, validation.When(m.CustomerEmail != "", is.EmailFormat.Error("客户邮箱地址格式错误"))),
		validation.Field(&m.CustomerPhone, validation.When(m.CustomerPhone != "", validation.Match(regexp.MustCompile(`^+\d{2}\d{11}$`)).Error("客户手机号码格式错误"))),
		validation.Field(&m.ShippingDate, validation.When(m.ShippingDate != "", validation.Date("2006-01-02 15:04").Error("包裹发货时间格式错误"))),
		validation.Field(&m.TrackingShippingDate, validation.When(m.TrackingShippingDate != "", validation.Date("20060102").Error("包裹发货时间格式错误"))),
	)
}

type Result struct {
	TrackingNumber string `json:"tracking_number"` //	包裹物流单号
	CourierCode    string `json:"courier_code"`    //	物流商对应的唯一简码
	OrderNumber    string `json:"order_number"`    //	包裹的订单号，由商家/平台所产生的订单编号
}

type CreateResult struct {
	Success []Result `json:"success"`
	Error   []Result `json:"error"`
}

func (s trackingService) Create(req CreateTrackingRequest) (res CreateResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().
		SetBody(req).
		Put("/create")
	if err != nil {
		return
	}

	r := struct {
		NormalResponse
		Data CreateResult `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &r); err == nil {
		res = r.Data
	}
	return
}

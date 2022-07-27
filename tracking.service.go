package tracking51

import (
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strings"
)

type trackingService service

type CreateTrackRequest struct {
	TrackingNumber          string `json:"tracking_number"`                     // 包裹物流单号
	CourierCode             string `json:"courier_code"`                        // 物流商对应的唯一简码
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

func (m CreateTrackRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
		validation.Field(&m.CourierCode, validation.Required.Error("物流商简码不能为空")),
		validation.Field(&m.CustomerEmail, validation.When(m.CustomerEmail != "", is.EmailFormat.Error("客户邮箱地址格式错误"))),
		validation.Field(&m.CustomerPhone, validation.When(m.CustomerPhone != "", validation.Match(regexp.MustCompile(`^+\d{2}\d{11}$`)).Error("客户手机号码格式错误"))),
		validation.Field(&m.ShippingDate, validation.When(m.ShippingDate != "", validation.Date("2006-01-02 15:04").Error("包裹发货时间格式错误"))),
		validation.Field(&m.TrackingShippingDate, validation.When(m.TrackingShippingDate != "", validation.Date("20060102").Error("跟踪包裹发货时间格式错误"))),
	)
}

type Result struct {
	TrackingNumber string `json:"tracking_number"` // 包裹物流单号
	CourierCode    string `json:"courier_code"`    // 物流商对应的唯一简码
	OrderNumber    string `json:"order_number"`    // 包裹的订单号，由商家/平台所产生的订单编号
}

type CreateResult struct {
	Success []Result `json:"success"`
	Error   []Result `json:"error"`
}

func (s trackingService) Create(req CreateTrackRequest) (res CreateResult, err error) {
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

// 获取查询结果
// https://www.51tracking.com/v3/api-index?language=Golang#%E8%8E%B7%E5%8F%96%E6%9F%A5%E8%AF%A2%E7%BB%93%E6%9E%9C

type Track struct {
	TrackingNumber     string          `json:"tracking_number"`     // 包裹物流单号
	CourierCode        string          `json:"courier_code"`        // 物流商对应的唯一简码
	LogisticsChannel   string          `json:"logistics_channel"`   // 自定义字段，用于填写物流渠道（比如某货代商）
	Destination        string          `json:"destination"`         // 目的国的二字简码
	TrackUpdate        bool            `json:"track_update"`        // 自动更新查询功能的状态，“true”代表系统会自动更新查询结果，“false”则反之
	Consignee          string          `json:"consignee"`           // 签收人
	Updating           bool            `json:"updating"`            // “true”表示该单号会被继续更新，“false”表示该单号已停止更新
	CreatedAt          string          `json:"created_at"`          // 创建查询的时间
	UpdateDate         string          `json:"update_date"`         // 系统最后更新查询的时间
	OrderCreateTime    string          `json:"order_create_time"`   // 包裹发货时间
	CustomerEmail      string          `json:"customer_email"`      // 客户邮箱
	CustomerPhone      string          `json:"customer_phone"`      // 顾客接收短信的手机号码
	Title              string          `json:"title"`               // 包裹名称
	OrderNumber        string          `json:"order_number"`        // 包裹的订单号，由商家/平台所产生的订单编号
	Note               string          `json:"note"`                // 备注，可自定义
	CustomerName       string          `json:"customer_name"`       // 客户姓名
	Archived           bool            `json:"archived"`            // “true”表示该单号已被归档，“false”表示该单号处于未归档状态
	Original           string          `json:"original"`            // 发件国的名称
	DestinationCountry string          `json:"destination_country"` // 目的国的名称
	TransitTime        int             `json:"transit_time"`        // 包裹的从被揽收至被送达的时长（天）
	StayTime           int             `json:"stay_time"`           // 物流信息未更新的时长（单位：天），由当前时间减去物流信息最近更新时间得到
	OriginInfo         TrackOriginInfo `json:"origin_info"`         // 发件国的物流信息
}

type TrackOriginInfo struct {
	DestinationTrackNumber string      `json:"destination_track_number"` // 该包裹在目的国的物流单号
	ReferenceNumber        string      `json:"reference_number"`         // 包裹对应的另一个单号，作用与当前单号相同（仅有少部分物流商提供）
	ExchangeNumber         string      `json:"exchangeNumber"`           // 该包裹在中转站的物流商单号
	ReceivedDate           string      `json:"received_date"`            // 物流商接收包裹的时间（也称为上网时间）
	DispatchedDate         string      `json:"dispatched_date"`          // 包裹封发时间，封发指将多个小包裹打包成一个货物（方便运输）
	DepartedAirportDate    string      `json:"departed_airport_date"`    // 包裹离开此出发机场的时间
	ArrivedAbroadDate      string      `json:"arrived_abroad_date"`      // 包裹达到目的国的时间
	CustomsReceivedDate    string      `json:"customs_received_date"`    // 包裹移交给海关的时间
	ArrivedDestinationDate string      `json:"arrived_destination_date"` // 包裹达到目的国、目的城市的时间
	Weblink                string      `json:"weblink"`                  // 物流商的官网的链接
	CourierPhone           string      `json:"courier_phone"`            // 物流商官网上的电话
	TrackInfo              []TrackInfo `json:"trackinfo"`                // 详细物流信息
	ServiceCode            string      `json:"service_code"`             // 快递服务类型，比如次日达（部分物流商返回）
	StatusInfo             string      `json:"status_info"`              // 最新的一条物流信息
	Weight                 string      `json:"weight"`                   // 该货物的重量（多个包裹会被打包成一个“货物”）
	DestinationInfo        string      `json:"destination_info"`         // 目的国的物流信息
	LatestEvent            string      `json:"latest_event"`             // 最新物流信息的梗概，包括以下信息：状态、地址、时间
	LatestCheckpointTime   string      `json:"lastest_checkpoint_time"`  // 最新物流信息的更新时间
}

// TrackInfo 详细物流信息
type TrackInfo struct {
	CheckpointDate              string `json:"checkpoint_date"`               // 本条物流信息的更新时间，由物流商提供（包裹被扫描时，物流信息会被更新）
	TrackingDetail              string `json:"tracking_detail"`               // 具体的物流情况
	Location                    string `json:"location"`                      // 物流信息更新的地址（该包裹被扫描时，所在的地址）
	CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`    // 根据具体物流情况所识别出来的物流状态
	CheckpointDeliverySubStatus string `json:"checkpoint_delivery_substatus"` // 物流状态的子状态（物流状态）
}

type TracksQueryParams struct {
	TrackingNumbers string `url:"tracking_numbers,omitempty"`  // 查询单号，每次不得超过40个，单号间以逗号分隔
	OrderNumbers    string `url:"order_numbers,omitempty"`     // 订单号，每次查询不得超过40个，订单号间以逗号分隔
	DeliveryStatus  string `url:"delivery_status,omitempty"`   // 发货状态
	ArchivedStatus  string `url:"archived_status,omitempty"`   // 指定该单号是否被归档。如果参数为字符串“true”，该单号将处于“归档”状态；如果参数为“false”，该单号处于“未归档”状态
	ItemsAmount     int    `url:"items_amount,omitempty"`      // 每页展示的单号个数
	PagesAmount     int    `url:"pages_amount,omitempty"`      // 返回结果的页数
	CreatedDateMin  int    `url:"created_date_min,omitempty"`  // 创建查询的起始时间，时间戳格式
	CreatedDateMax  int    `url:"created_date_max,omitempty"`  // 创建查询的结束时间，时间戳格式
	ShippingDateMin int    `url:"shipping_date_min,omitempty"` // 发货的起始时间，时间戳格式
	ShippingDateMax int    `url:"shipping_date_max,omitempty"` // 发货的结束时间，时间戳格式
	UpdatedDateMin  int    `url:"updated_date_min,omitempty"`  // 查询更新的起始时间，时间戳格式
	UpdatedDateMax  int    `url:"updated_date_max,omitempty"`  // 查询更新的结束时间，时间戳格式
	Lang            string `url:"lang,omitempty"`              // 查询结果的语言（例子：cn, en），若未指定该参数，结果会以英文或中文呈现。 注意：只有物流商支持多语言查询结果时，该指定才会生效
}

func (m TracksQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumbers, validation.When(m.TrackingNumbers != "", validation.By(func(value interface{}) error {
			numbers, ok := value.(string)
			if !ok {
				return fmt.Errorf("无效的查询单号：%s", m.TrackingNumbers)
			}
			if len(strings.Split(numbers, ",")) > 40 {
				return errors.New("查询单号不能超过 40 个")
			}
			return nil
		}))),
		validation.Field(&m.OrderNumbers, validation.When(m.OrderNumbers != "", validation.By(func(value interface{}) error {
			numbers, ok := value.(string)
			if !ok {
				return fmt.Errorf("无效的订单号：%s", m.OrderNumbers)
			}
			if len(strings.Split(numbers, ",")) > 40 {
				return errors.New("订单号不能超过 40 个")
			}
			return nil
		}))),
		validation.Field(&m.DeliveryStatus, validation.When(m.DeliveryStatus != "", validation.In(StatusPending, StatusNotFound, StatusTransit, StatusPickup, StatusDelivered, StatusExpired, StatusUndelivered, StatusException, StatusInfoReceived).Error("无效的发货状态"))),
		validation.Field(&m.ArchivedStatus, validation.When(m.ArchivedStatus != "", validation.In("true", "false").Error("无效的归档状态"))),
		validation.Field(&m.Lang, validation.When(m.Lang != "", validation.In(ChineseLanguage, EnglishLanguage).Error("无效的查询结果语言"))),
	)
}

func (s trackingService) Query(params TracksQueryParams) (items []Track, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	if params.PagesAmount <= 0 {
		params.PagesAmount = 1
	}
	if params.ItemsAmount <= 0 {
		params.ItemsAmount = 100
	}
	resp, err := s.httpClient.R().
		SetQueryParamsFromValues(toValues(params)).
		Get("/get")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data []Track `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		items = res.Data
		isLastPage = len(items) < params.ItemsAmount
	}
	return
}

type trackingNumberCourierCode struct {
	TrackingNumber string `json:"tracking_number"` // 包裹物流单号
	CourierCode    string `json:"courier_code"`    // 物流商对应的唯一简码
}

// 删除查询单号

type DeleteTrackRequest trackingNumberCourierCode

type DeleteTrackRequests []DeleteTrackRequest

func (m DeleteTrackRequests) Validate() error {
	n := len(m)
	if n == 0 {
		return errors.New("请求数据不能为空")
	} else if n > 40 {
		return errors.New("请求数据不能超过 40 个")
	}

	var err error
	for _, request := range m {
		err = validation.ValidateStruct(&request,
			validation.Field(&request.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
			validation.Field(&request.CourierCode, validation.Required.Error("物流商简码不能为空")),
		)
		if err != nil {
			break
		}
	}
	return err
}

type DeleteTrackResult trackingNumberCourierCode

func (s trackingService) Delete(req DeleteTrackRequests) (success []DeleteTrackResult, error []DeleteTrackResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().
		SetBody(req).
		Delete("/delete")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data struct {
			Success []DeleteTrackResult `json:"success"`
			Error   []DeleteTrackResult `json:"error"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		success = res.Data.Success
		error = res.Data.Error
	}
	return
}

// 停止更新

type StopUpdateTrackRequest trackingNumberCourierCode
type StopUpdateTrackRequests []StopUpdateTrackRequest

func (m StopUpdateTrackRequests) Validate() error {
	n := len(m)
	if n == 0 {
		return errors.New("请求数据不能为空")
	} else if n > 40 {
		return errors.New("请求数据不能超过 40 个")
	}

	var err error
	for _, request := range m {
		err = validation.ValidateStruct(&request,
			validation.Field(&request.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
			validation.Field(&request.CourierCode, validation.Required.Error("物流商简码不能为空")),
		)
		if err != nil {
			break
		}
	}
	return err
}

type StopUpdateResultSuccess trackingNumberCourierCode
type StopUpdateResultError trackingNumberCourierCode

func (s trackingService) StopUpdate(req StopUpdateTrackRequests) (success []StopUpdateResultSuccess, error []StopUpdateResultError, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().SetBody(req).Post("/notupdate")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data struct {
			Success []StopUpdateResultSuccess `json:"success"`
			Error   []StopUpdateResultError   `json:"error"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		success = res.Data.Success
		error = res.Data.Error
	}
	return
}

// 手动更新

type RefreshTrackRequest trackingNumberCourierCode

func (m RefreshTrackRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
		validation.Field(&m.CourierCode, validation.Required.Error("物流商简码不能为空")),
	)
}

type RefreshTrackRequests []RefreshTrackRequest

func (m RefreshTrackRequests) Validate() error {
	n := len(m)
	if n == 0 {
		return errors.New("请求数据不能为空")
	} else if n > 40 {
		return errors.New("请求数据不能超过 40 个")
	}

	var err error
	for _, request := range m {
		err = validation.ValidateStruct(&request,
			validation.Field(&request.TrackingNumber, validation.Required.Error("包裹物流单号不能为空")),
			validation.Field(&request.CourierCode, validation.Required.Error("物流商简码不能为空")),
		)
		if err != nil {
			break
		}
	}
	return err
}

type RefreshResultSuccess struct {
	trackingNumberCourierCode
}

type RefreshResultError struct {
	trackingNumberCourierCode
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (s trackingService) Refresh(req RefreshTrackRequests) (success []RefreshResultSuccess, error []RefreshResultError, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().SetBody(req).Post("/manualupdate")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data struct {
			Success []RefreshResultSuccess `json:"success"`
			Error   []RefreshResultError   `json:"error"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		success = res.Data.Success
		error = res.Data.Error
	}
	return
}

// 统计包裹状态

type StatusStatisticRequest struct {
	CourierCode     string `url:"courier_code,omitempty"`      // 物流商对应的唯一简码
	CreatedDateMin  int    `url:"created_date_min,omitempty"`  // 创建查询的起始时间，时间戳格式
	CreatedDateMax  int    `url:"created_date_max,omitempty"`  // 创建查询的结束时间，时间戳格式
	ShippingDateMin int    `url:"shipping_date_min,omitempty"` // 发货的起始时间，时间戳格式
	ShippingDateMax int    `url:"shipping_date_max,omitempty"` // 发货的结束时间，时间戳格式
}

type StatusStatistic struct {
	Pending      int `json:"pending"`      // 查询中：新增包裹正在查询中，请等待
	NotFound     int `json:"notfound"`     // 查询不到：包裹信息目前查询不到
	Transit      int `json:"transit"`      // 运输途中：物流商已揽件，包裹正被发往目的地
	Pickup       int `json:"pickup"`       // 到达待取：包裹正在派送中，或到达当地收发点
	Delivered    int `json:"delivered"`    // 成功签收：包裹已被成功投递
	Expired      int `json:"expired"`      // 运输过久：包裹在很长时间内都未投递成功。快递包裹超过30天、邮政包裹超过60天未投递成功，该查询会被识别为此状态
	Undelivered  int `json:"undelivered"`  // 投递失败：快递员投递失败（通常会留有通知并再次尝试投递）
	Exception    int `json:"exception"`    // 可能异常：包裹退回、包裹丢失、清关失败等异常情况
	InfoReceived int `json:"infoReceived"` // 待上网：包裹正在等待被揽件
}

func (m StatusStatisticRequest) Validate() error {
	return nil
}

func (s trackingService) StatusStatistic(req StatusStatisticRequest) (stat StatusStatistic, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().SetBody(req).Get("/status")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data StatusStatistic `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		stat = res.Data
	}
	return
}

// 时效

type TransitTime struct {
	CourierCode         string  `json:"courier_code"`          // 物流商对应的唯一简码
	OriginalCode        string  `json:"original_code"`         // 发件国二字简码
	DestinationCode     string  `json:"destination_code"`      // 目的国二字简码
	Total               int     `json:"total"`                 // 未签收的总单号数量
	Delivered           int     `json:"delivered"`             // 已签收的总单号数量
	Range1To7           float64 `json:"range_1_7"`             // 送达时间为0～7天的单号的占比
	Range8To15          float64 `json:"range_8_15"`            // 送达时间为7～15天的单号的占比
	Range16To30         float64 `json:"range_16_30"`           // 送达时间为16～30天的单号的占比
	Range31To60         float64 `json:"range_31_60"`           // 送达时间为31～60天的单号的占比
	Range60Up           float64 `json:"range_60_up"`           // 送达时间为31～60天的单号的占比
	AverageDeliveryTime float64 `json:"average_delivery_time"` // 平均送达时间（单位：天）
}

type TransitTimeRequest struct {
	CourierCode     string `json:"courier_code"`      // 物流商对应的唯一简码
	OriginalCode    string `json:"original_code"`     // 发件国二字简码
	DestinationCode string `json:"destination_code "` // 目的国的二字简码
}

type TransitTimeRequests []TransitTimeRequest

func (m TransitTimeRequests) Validate() error {
	if len(m) == 0 {
		return errors.New("请求数据不能为空")
	}

	var err error
	for _, request := range m {
		err = validation.ValidateStruct(&request,
			validation.Field(&request.CourierCode, validation.Required.Error("物流商简码不能为空")),
			validation.Field(&request.OriginalCode, validation.Required.Error("发件国二字简码不能为空")),
			validation.Field(&request.DestinationCode, validation.Required.Error("目的国二字简码不能为空")),
		)
		if err != nil {
			break
		}
	}
	return err
}

func (s trackingService) TransitTime(req TransitTimeRequests) (success []TransitTime, error []TransitTime, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	resp, err := s.httpClient.R().
		SetBody(req).
		Get("/transittime")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data struct {
			Success []TransitTime `json:"success"`
			Error   []TransitTime `json:"error"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		success = res.Data.Success
		error = res.Data.Error
	}
	return
}

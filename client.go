package tracking51

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"github.com/hiscaler/51tracking-go/config"
	"github.com/hiscaler/gox/bytex"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// https://www.51tracking.com/v3/api-index?language=Golang#%E5%93%8D%E5%BA%94
const (
	Success                             = 200 // 无错误
	PaymentRequiredError                = 203 // API 服务只提供给付费账户，请付费购买单号以解锁 API 服务
	NoContent                           = 204 // 请求成功，但未获取到数据，可能是该单号、所查询目标数据不存在
	BadRequestError                     = 400 // 请求类型错误
	UnauthorizedError                   = 401 // 授权失败或没有权限，请检查并确保你 API Key 正确无误
	NotFoundError                       = 404 // 请求的资源不存在
	TimeOutError                        = 408 // 请求超时
	RequestParametersTooLongError       = 411 // 请求参数长度超过限制
	RequestParametersFormatError        = 412 // 请求参数格式不合要求
	RequestParametersExceededLimitError = 413 // 请求参数数量超过限制
	LostRequestParametersOrParseError   = 417 // 缺少请求参数或者请求参数无法解析
	ParametersInvalidError              = 421 // 部分必填参数为空
	CourierCodeInvalidError             = 422 // 物流商简码无法识别或者不支持该物流商
	TrackingNumberIsExistsError         = 423 // 跟踪单号已存在，无需再次创建
	TrackingNumberIsNotExistsError      = 424 // 跟踪单号不存在
	TooManyRequestsError                = 429 // API 请求频率次数限制，请稍后再试
	InternalError                       = 511 // 系统错误
)

const (
	StatusPending      = "pending"      // 查询中
	StatusNotFound     = "notfound"     // 查询不到
	StatusTransit      = "transit"      // 运输中
	StatusPickup       = "pickup"       // 到达待取
	StatusDelivered    = "delivered"    // 成功签收
	StatusExpired      = "notfound"     // 运输过久
	StatusUndelivered  = "undelivered"  // 投递失败
	StatusException    = "exception"    // 可能异常
	StatusInfoReceived = "inforeceived" // 待上网
)

const (
	ChineseLanguage = "cn"
	EnglishLanguage = "en"
)

const (
	Version   = "0.0.1"
	userAgent = "51tracking API Client-Golang/" + Version + " (https://github.com/hiscaler/51tracking-go)"
)

type Tracking51 struct {
	latestRequestTime time.Time      // 最后请求时间（在传入了 IntervalTime 后，该值用于控制接口调取频率处理，未传入的话不起作用）
	config            *config.Config // 配置
	httpClient        *resty.Client  // Resty Client
	Services          services       // API Services
}

func NewTracking51(config config.Config) *Tracking51 {
	logger := log.New(os.Stdout, "[ 51Tracking ] ", log.LstdFlags|log.Llongfile)
	client := &Tracking51{
		config: &config,
	}

	baseURL := "https://api.51tracking.com/v3/trackings"
	if config.Sandbox {
		baseURL += "/sandbox"
	}
	httpClient := resty.New().
		SetDebug(config.Debug).
		SetBaseURL(baseURL).
		SetHeaders(map[string]string{
			"Content-Type":     "application/json",
			"Accept":           "application/json",
			"User-Agent":       userAgent,
			"Tracking-Api-Key": config.AppKey,
		}).
		SetTimeout(10 * time.Second).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) (err error) {
			if response.IsError() {
				return fmt.Errorf("%s: %s", response.Status(), bytex.ToString(response.Body()))
			}

			r := struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{}
			if err = json.Unmarshal(response.Body(), &r); err == nil {
				if r.Code != Success {
					err = ErrorWrap(r.Code, r.Message)
				}
			} else {
				logger.Printf("JSON Unmarshal error: %s", err.Error())
			}

			if err != nil {
				logger.Printf("OnAfterResponse error: %s", err.Error())
			}
			return
		}).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return false
			}

			retry := response.StatusCode() == http.StatusTooManyRequests
			if !retry {
				r := struct{ Code int }{}
				retry = json.Unmarshal(response.Body(), &r) == nil && r.Code == TooManyRequestsError
			}
			if retry {
				text := response.Request.URL
				if err != nil {
					text += fmt.Sprintf(", error: %s", err.Error())
				}
				logger.Printf("Retry request: %s", text)
			}
			return retry
		})

	if config.IntervalTime > 0 {
		httpClient.OnBeforeRequest(func(c *resty.Client, request *resty.Request) error {
			now := time.Now()
			if client.latestRequestTime.IsZero() {
				client.latestRequestTime = now
				return nil
			}

			if config.Debug {
				logger.Printf("URL: %s", request.URL)
				logger.Printf("Client latest request time: %s, Current request time: %s", client.latestRequestTime.Format("2006-01-02 15:04:05.000"), now.Format("2006-01-02 15:04:05.000"))
			}
			d := now.Sub(client.latestRequestTime)
			if d.Milliseconds() < config.IntervalTime {
				d = time.Duration(config.IntervalTime-d.Milliseconds()) * time.Millisecond
				if config.Debug {
					logger.Printf("Sleep %d milliseconds", d.Milliseconds())
				}
				time.Sleep(d)
			}
			client.latestRequestTime = now
			return nil
		})
	}
	client.httpClient = httpClient
	xService := service{
		config:     &config,
		logger:     logger,
		httpClient: client.httpClient,
	}
	client.Services = services{
		Account:  (accountService)(xService),
		Courier:  (courierService)(xService),
		Tracking: (trackingService)(xService),
	}
	return client
}

// SetDebug 设置是否开启调试模式
func (t *Tracking51) SetDebug(v bool) *Tracking51 {
	t.config.Debug = v
	t.httpClient.SetDebug(v)
	return t
}

type NormalResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorWrap 错误包装
func ErrorWrap(code int, message string) error {
	if code == Success || code == NoContent {
		return nil
	}

	switch code {
	case PaymentRequiredError:
		message = "API 服务只提供给付费账户，请付费购买单号以解锁 API 服务"
	case BadRequestError:
		message = "请求类型错误"
	case UnauthorizedError:
		message = "授权失败或没有权限，请检查并确保你 API Key 正确无误"
	case NotFoundError:
		message = "请求的资源不存在"
	case TimeOutError:
		message = "请求超时"
	case RequestParametersTooLongError:
		message = "请求参数长度超过限制"
	case RequestParametersFormatError:
		message = "请求参数格式不合要求"
	case RequestParametersExceededLimitError:
		message = "请求参数数量超过限制"
	case LostRequestParametersOrParseError:
		message = "缺少请求参数或者请求参数无法解析"
	case ParametersInvalidError:
		message = "部分必填参数为空"
	case CourierCodeInvalidError:
		message = "物流商简码无法识别或者不支持该物流商"
	case TrackingNumberIsExistsError:
		message = "跟踪单号已存在，无需再次创建"
	case TrackingNumberIsNotExistsError:
		message = "跟踪单号不存在"
	case TooManyRequestsError:
		message = "API 请求频率次数限制，请稍后再试"
	case InternalError:
		message = "系统错误"
	}
	return fmt.Errorf("%d: %s", code, message)
}

// change to url.values
func toValues(i interface{}) (values url.Values) {
	values, _ = query.Values(i)
	return
}

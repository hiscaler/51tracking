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
	OK                                  = 200 // 无错误
	PaymentRequiredError                = 203 // API 服务只提供给付费账户，请付费购买单号以解锁 API 服务
	NoContentError                      = 204 // 请求成功，但未获取到数据，可能是该单号、所查询目标数据不存在
	BadRequestError                     = 400 // 请求类型错误
	UnauthorizedError                   = 401 // 授权失败或没有权限，请检查并确保你 API Key 正确无误
	NotFoundError                       = 404 // 该页面不存在
	TimeOutError                        = 408 // 请求超时
	RequestParametersTooLongError       = 411 // 请求参数长度超过限制
	RequestParametersFormatError        = 412 // 请求参数格式不合要求
	RequestParametersExceededLimitError = 413 // 请求参数数量超过限制
	TooManyRequestsError                = 429 // API请求频率次限制，请稍后再试
)

const (
	Version   = "0.0.1"
	userAgent = "51tracking API Client-Golang/" + Version + " (https://github.com/hiscaler/51tracking-go)"
)

type Tracking51 struct {
	config     *config.Config // 配置
	httpClient *resty.Client  // Resty Client
	Services   services       // API Services
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
				if r.Code != OK {
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
		SetRetryWaitTime(5 * time.Second).
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
	if code == OK {
		return nil
	}

	switch code {
	case TooManyRequestsError:
		message = "接口请求超请求次数限额"
	}
	return fmt.Errorf("%d: %s", code, message)
}

// change to url.values
func toValues(i interface{}) (values url.Values) {
	values, _ = query.Values(i)
	return
}

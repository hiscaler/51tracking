51Tracking SDK
====================

## 文档地址

https://www.51tracking.com/v3/api-index?language=Golang#api-version

## 数据更新频率

### 正常

- 添加 30 天内未签收的属于正常更新频次也就是下限是 4 小时上限是 6 小时；
- 30~45 天仍未签收就会降低更新频率至下限 6 小时上限 8 小时；
- 45~60 天仍未签收更新频率变为下限 12 小时上限 14 小时；
- 60~80 天仍未签收更新频率变为下限 24 小时上限 26 小时；
- 超过 80 天仍未签收将停止更新。

### 其他

- 快递更新间隔正常时间是 4~6 个小时；
- 邮政更新间隔是 6~8 个小时；
- sfb2c 更新间隔为 8 小时。

## 手动更新条件

1. 查询不到超过 15 天，已经停止更新了的单号；
2. 已经签收，并且已经停止更新了的单号；
3. 添加时间超过 30 天，已经停止更新的单号；
4. 其他条件暂时无法通过该接口更新。

## 安装

```go
go get github.com/hiscaler/51tracking-go
```

## 使用

```
client := NewTracking51(config.Config{
    Debug:        true,
    Sandbox:      true,
    AppKey:       "xxx",
    IntervalTime: 1500,
})
```

## 配置说明

```go
type Config struct {
	Debug        bool   // 是否为调试模式（调试模式下会输出 HTTP 请求和返回数据）
	Sandbox      bool   // 是否为沙箱测试环境
	Version      string // API 版本（当前固定为 V3）
	AppKey       string // App Key
	IntervalTime int64  // 当前请求与上次请求间隔的时间（单位为毫秒），默认为零，表示没有间隔，大于 0 表示实际间隔的毫秒数
}
```

## 服务

### Account

- 帐号情况

```go
client.Services.Account.Profile()
```

### Courier

- 获取物流商列表

```go
client.Services.Courier.List()
```

- 修改包裹物流商

```go
client.Services.Courier.Change("trackingNumber", "oldCourierCode", "newCourierCode")
```

### Tracking

- 添加物流单号

```go
client.Services.Tracking.Create()
```

- 修改单号信息

```go
client.Services.Tracking.Update()
```

- 获取查询结果

```go
params := TracksQueryParams{}
client.Services.Tracking.Query(params)
```

- 删除查询单号

```go
client.Services.Tracking.Delete([]DeleteTrackRequest{})
```

- 停止单号更新

```go
client.Services.Tracking.StopUpdate(StopUpdateRequests{})
```

- 手动更新

```go
client.Services.Tracking.Refresh([]RefreshRequest{})
```

- 统计包裹状态

```go
client.Services.Tracking.StatusStatistic(StatusStatisticRequest{})
```

- 时效

```go
client.Services.Tracking.TransitTime(TransitTimeRequest{})
```

- 检测偏远地区

```go
client.Services.Tracking.RemoteDetection(RemoteDetectionRequest{})
```

## Webhook

针对 51Tracking 的数据推送，提供了 WebhookRequest 结构体，您可以使用他来接受推送过来的数据，并判断 Code 是否为 200 且 Data.Valid() 是否有效来进行下一步的业务逻辑处理。

Data.Valid(you51TrackingAccountEmail) 用来判断推送的数据是否来自 51Tracking，根据需要，您可以跳过此步。

```go
var wr WebRequest
if json.Unmarshal(resp.Body, &wr) != nil {
	if wr.Code == 200 && wr.Data.Valid(you51TrackingAccountEmail) {
	    // you code	
    }
}
```
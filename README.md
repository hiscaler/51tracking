51Tracking SDK
====================

## 文档地址

https://www.51tracking.com/v3/api-index?language=Golang#api-version

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

### 帐号

- 帐号情况

```go
client.Services.Account.Profile()
```

### 物流商

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

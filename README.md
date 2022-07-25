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
    Debug:   true,
    Sandbox: true,
    AppKey:  "xxx",
})
```

## 服务

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
client.Services.Tracking.All(params)
```

- 删除查询单号

```go
client.Services.Tracking.Delete([]DeleteTrackRequest{})
```

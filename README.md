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
client := NewTracking51(c)
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

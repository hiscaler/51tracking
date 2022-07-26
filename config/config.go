package config

type Config struct {
	Debug        bool   // 是否为调试模式（调试模式下会输出 HTTP 请求和返回数据）
	Sandbox      bool   // 是否为沙箱测试环境
	Version      string // API 版本（当前固定为 V3）
	AppKey       string // App Key
	IntervalTime int64  // 当前请求与上次请求间隔的时间（单位为毫秒），默认为零，表示没有间隔，大于 0 表示实际间隔的毫秒数（大于零小于 1000 的值强制设置为 1000）
}

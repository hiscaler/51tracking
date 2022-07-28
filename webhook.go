package tracking51

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hiscaler/gox/stringx"
	"strconv"
)

// Webhook 数据处理

type webhookVerify struct {
	Timestamp int    `json:"timestamp"`
	Signature string `json:"signature"`
	UserTag   string `json:"usertag"`
}

type Webhook struct {
	Track
	Verify webhookVerify `json:"verify"`
}

type WebhookRequest struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Webhook `json:"data"`
}

// Valid 传递 51tracking 用户邮箱验证是否为有效的 51Tracking 推送
func (wr Webhook) Valid(email string) bool {
	if wr.Verify.Timestamp == 0 || wr.Verify.Signature == "" {
		return false
	}
	hash := hmac.New(sha256.New, stringx.ToBytes(email))
	hash.Write(stringx.ToBytes(strconv.Itoa(wr.Verify.Timestamp)))
	return hex.EncodeToString(hash.Sum(nil)) == wr.Verify.Signature
}

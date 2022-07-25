package tracking51

import "encoding/json"

type accountService service

type AccountProfile struct {
	Email       string `json:"email"`        // 登录邮箱
	RegTime     int    `json:"regtime"`      // 账户注册时间
	Phone       string `json:"phone"`        // 账户绑定的手机号码
	SMS         int    `json:"sms"`          // 短信剩余条数
	TrackNumber int    `json:"track_number"` // 账户剩余的单号额度
}

func (s accountService) Profile() (profile AccountProfile, err error) {
	resp, err := s.httpClient.R().
		Get("/userinfo")
	if err != nil {
		return
	}

	res := struct {
		NormalResponse
		Data AccountProfile `json:"data"`
	}{}
	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		profile = res.Data
	}
	return
}

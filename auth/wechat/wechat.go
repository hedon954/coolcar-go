package wechat

import (
	"fmt"

	"github.com/medivhzhan/weapp/v2"
)

type Service struct {
	AppID     string
	AppSecret string
}

/**
GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

appid
secret
js_code
grant_type=authorization_code
*/

// Resolve 跳动 auth.code2Session 解析 OpenID
func (s *Service) Resolve(code string) (string, error) {
	resp, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "", fmt.Errorf("weapp.Login error: %v", err)
	}
	return resp.OpenID, nil
}

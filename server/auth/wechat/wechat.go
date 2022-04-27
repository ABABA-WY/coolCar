package wechat

import (
	"fmt"
	"github.com/medivhzhan/weapp/v2"
)

type Service struct {
	AppID     string
	AppSecret string
}

func (s *Service) Resolver(code string) (string, error) {
	response, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "", fmt.Errorf("weapp login error:%v", err)
	}
	if response.GetResponseError() != nil {
		return "", fmt.Errorf("weapp response error:%v", response.GetResponseError())
	}
	return response.OpenID, nil
}

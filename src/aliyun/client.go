package aliyun

import (
	"fmt"
	cas "github.com/alibabacloud-go/cas-20200407/v3/client"
	cdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var ErrCertExists = fmt.Errorf("cert exists")

var casClient *cas.Client
var cdnClient *cdn.Client

func createCASClient(key string, secret string) (*cas.Client, error) {
	result, err := cas.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(key),
		AccessKeySecret: tea.String(secret),
		Endpoint:        tea.String("cas.aliyuncs.com"),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func createCDNClient(key string, secret string) (*cdn.Client, error) {
	result, err := cdn.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(key),
		AccessKeySecret: tea.String(secret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

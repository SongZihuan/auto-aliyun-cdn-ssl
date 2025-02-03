package aliyun

import (
	"fmt"
	cas "github.com/alibabacloud-go/cas-20200407/v3/client"
	cdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

var ErrCertExists = fmt.Errorf("cert exists")

var casClient *cas.Client
var cdnClient *cdn.Client
var dcdnClient *dcdn.Client

func createClient(key string, secret string) (err error) {
	casClient, err = createCASClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CAS client error: %s", err.Error())
	}

	cdnClient, err = createCDNClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CDN client error: %s", err.Error())
	}

	dcdnClient, err = createDCDNClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CDN client error: %s", err.Error())
	}

	return nil
}

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

func createDCDNClient(key string, secret string) (*dcdn.Client, error) {
	result, err := dcdn.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(key),
		AccessKeySecret: tea.String(secret),
		Endpoint:        tea.String("dcdn.aliyuncs.com"),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

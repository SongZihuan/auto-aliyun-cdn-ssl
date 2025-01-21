package aliyun

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"os"
)

var international = false

func Init() (err error) {
	if !config.IsReady() {
		panic("config is not ready")
	}

	international = config.GetConfig().Aliyun.International.ToBool(false)
	key := config.GetConfig().Aliyun.Key
	secret := config.GetConfig().Aliyun.Secret

	casClient, err = createCASClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CAS client error: %s\n", err.Error())
	}

	cdnClient, err = createCDNClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CDN client error: %s\n", err.Error())
	}

	return nil
}

func UpdateCDNHttpsByFilePath(domain string, cert string, prikey string) error {
	certData, err := os.ReadFile(cert)
	if err != nil {
		return fmt.Errorf("read cert file error: %s\n", err.Error())
	}

	privateKeyData, err := os.ReadFile(prikey)
	if err != nil {
		return fmt.Errorf("read private key error: %s\n", err.Error())
	}

	return UpdateCDNHttps(domain, certData, privateKeyData)
}

func UpdateCDNHttps(domain string, certData []byte, privateKeyData []byte) error {
	certID, certName, err := uploadCert(casClient, certData, privateKeyData)
	if err != nil && errors.Is(err, ErrCertExists) {
		logger.Warn("证书已存在, 不在重新更新CDN")
		return nil
	} else if err != nil {
		return fmt.Errorf("upload error: %s\n", err.Error())
	}

	err = setDomainServerCertificate(cdnClient, domain, certID, certName)
	if err != nil {
		return err
	}

	return nil
}

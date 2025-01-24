package aliyun

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/database"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"os"
	"strings"
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
		return fmt.Errorf("init alibaba cloud sdk CAS client error: %s", err.Error())
	}

	cdnClient, err = createCDNClient(key, secret)
	if err != nil {
		return fmt.Errorf("init alibaba cloud sdk CDN client error: %s", err.Error())
	}

	return nil
}

func UpdateCDNHttpsByFilePath(domainList []string, cert string, prikey string) error {
	certData, err := os.ReadFile(cert)
	if err != nil {
		return fmt.Errorf("read cert file error: %s", err.Error())
	}

	privateKeyData, err := os.ReadFile(prikey)
	if err != nil {
		return fmt.Errorf("read private key error: %s", err.Error())
	}

	return UpdateCDNHttps(domainList, certData, privateKeyData)
}

func UpdateCDNHttps(domainList []string, certData []byte, privateKeyData []byte) error {
	certID, certName, subject, err := uploadCert(certData, privateKeyData)
	if err != nil && errors.Is(err, ErrCertExists) && certName != "" {
		logger.Warnf("证书已存在, 尝试检测CDN证书记录并更新（%s）", strings.Join(domainList, ", "))

		for _, domain := range domainList {
			cert, need, err := database.CheckNeedUpdateDomain(certName, domain)
			if err != nil {
				logger.Errorf("check aliyun cloud cdn domain ssl status from sqlite error: %s", err.Error())
			} else if need && cert != nil {
				setDomainServerCertificateNotError(domain, cert.CertID, cert.Name)
				err = database.UpdateDomain(cert.CertID, cert.Name, cert.Subject, domain)
				if err != nil {
					logger.Error("aliyun cloud ssl domain save to sqlite error: %s", err.Error())
				}
			} else {
				// 无需更新
			}
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("aliyun cloud ssl cert/key upload error: %s", err.Error())
	} else {
		err = database.UpdateCert(certID, certName, subject)
		if err != nil {
			logger.Errorf("aliyun cloud ssl cert/key save to sqlite error: %s", err.Error())
		}

		for _, domain := range domainList {
			setDomainServerCertificateNotError(domain, certID, certName)
			err = database.UpdateDomain(certID, certName, subject, domain)
			if err != nil {
				logger.Error("aliyun cloud ssl domain save to sqlite error: %s", err.Error())
			}
		}
	}

	return nil
}

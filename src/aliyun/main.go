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

	logger.Infof("成功从 %s 读取证书，从 %s 读取密钥，这些证书将被用于域名：%s", cert, prikey, strings.Join(domainList, ", "))
	return UpdateCDNHttps(domainList, certData, privateKeyData)
}

func UpdateCDNHttps(domainList []string, certData []byte, privateKeyData []byte) error {
	certID, certName, subject, err := uploadCert(certData, privateKeyData)
	if err != nil && errors.Is(err, ErrCertExists) && certName != "" {
		logger.Infof("证书已存在, 尝试检测 CDN域名 (%s) 证书记录并更新", strings.Join(domainList, ", "))

		for _, domain := range domainList {
			cert, need, err := database.CheckNeedUpdateDomain(certName, domain)
			if err != nil {
				logger.Errorf("在检测 域名 (%s) 是否应该更新时遇到了错误，但不影响后续检查: %s", domain, err.Error())
			} else if need && cert != nil {
				logger.Infof("确认域名 (%s) 需要更新， 证书Subject: %s, 证书名字：%s, 证书ID：%d", domain, cert.Subject, cert.Name, cert.CertID)
				setDomainServerCertificateNotError(domain, cert.CertID, cert.Name)
				err = database.UpdateDomain(cert.CertID, cert.Name, cert.Subject, domain)
				if err != nil {
					logger.Errorf("在更新 域名 (%s) 状态到数据库时遇到了错误，但不影响后续检查: %s", domain, err.Error())
				}
			} else if !need && cert != nil {
				// 无需更新
				logger.Infof("确认域名 (%s) 无需更新证书，并找到证书，证书Subject: %s, 证书名字：%s, 证书ID：%d", domain, cert.Subject, cert.Name, cert.CertID)
			} else if !need { // cert == nil
				// 无需更新
				logger.Infof("检测到 域名 (%s) 无需更新证书，未能找到证书相关记录", domain)
			}
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("aliyun cloud ssl cert/key upload error: %s", err.Error())
	} else {
		dbcerterr := database.UpdateCert(certID, certName, subject)
		if dbcerterr != nil {
			logger.Errorf("保存域名（%s）证书信息到数据库时发生了错误: %s", strings.Join(domainList, ", "), dbcerterr.Error())
		}

		for _, domain := range domainList {
			setDomainServerCertificateNotError(domain, certID, certName)
			if dbcerterr != nil {
				logger.Error("因为证书信息未能写入数据库，所以 域名 (%s) 信息不尝试写入数据库")
			} else {
				err = database.UpdateDomain(certID, certName, subject, domain)
				if err != nil {
					logger.Errorf("保存 域名 (%s) 信息到数据库发生了错误，但不影响后续域名更新: %s", dbcerterr, err.Error())
				}
			}
		}
	}

	return nil
}

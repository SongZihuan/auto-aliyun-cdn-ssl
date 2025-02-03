package aliyun

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/database"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"os"
)

var international = false

func Init() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	international = config.GetConfig().Aliyun.International.ToBool(false)
	key := config.GetConfig().Aliyun.Key
	secret := config.GetConfig().Aliyun.Secret

	err := createClient(key, secret)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDomainHttpsByFilePath(collection *config.DomainListCollection, cert string, prikey string) error {
	certData, err := os.ReadFile(cert)
	if err != nil {
		return fmt.Errorf("read cert file error: %s", err.Error())
	}

	privateKeyData, err := os.ReadFile(prikey)
	if err != nil {
		return fmt.Errorf("read private key error: %s", err.Error())
	}

	logger.Infof("成功从 %s 读取证书，从 %s 读取密钥，这些证书将被用于域名：%s", cert, prikey, collection.Domain2Str())
	return UpdateDomainHttps(collection, certData, privateKeyData)
}

func UpdateDomainHttps(collection *config.DomainListCollection, certData []byte, privateKeyData []byte) error {
	certID, certName, subject, err := uploadCert(certData, privateKeyData)
	if err != nil && errors.Is(err, ErrCertExists) && certName != "" {
		logger.Infof("证书已存在, 尝试检测 CDN域名 (%s) 证书记录并更新", collection.Domain2Str())

		for _, domain := range collection.Domain {
			if domain.Type == config.DomainTypeCDN {
				cert, need, err := database.CheckNeedUpdateCDNDomain(certName, domain.Domain)
				if err != nil {
					logger.Errorf("在检测 域名 (%s) 是否应该更新时遇到了错误，但不影响后续检查: %s", domain.Domain, err.Error())
				} else if cert == nil {
					logger.Infof("检测到 域名 (%s) 无需更新证书，未能找到证书相关记录", domain.Domain)
				} else if need {
					logger.Infof("确认域名 (%s) 需要更新， 证书Subject: %s, 证书名字：%s, 证书ID：%d", domain.Domain, cert.Subject, cert.Name, cert.CertID)
					err = setCDNServerCertificateNotPanic(domain.Domain, cert.CertID, cert.Name)
					if err == nil {
						err = database.UpdateCDNDomain(cert.CertID, cert.Name, cert.Subject, domain.Domain)
						if err != nil {
							logger.Errorf("在更新 域名 (%s) 状态到数据库时遇到了错误，但不影响后续检查: %s", domain.Domain, err.Error())
						}
					}
				} else {
					// 无需更新
					logger.Infof("确认域名 (%s) 无需更新证书，并找到证书，证书Subject: %s, 证书名字：%s, 证书ID：%d", domain.Domain, cert.Subject, cert.Name, cert.CertID)
				}
			} else if domain.Type == config.DomainTypeDCDN {
				cert, need, err := database.CheckNeedUpdateDCDNDomain(certName, domain.Domain)
				if err != nil {
					logger.Errorf("在检测 域名 (%s) 是否应该更新时遇到了错误，但不影响后续检查: %s", domain.Domain, err.Error())
				} else if cert == nil {
					logger.Infof("检测到 域名 (%s) 无需更新证书，未能找到证书相关记录", domain.Domain)
				} else if need {
					logger.Infof("确认域名 (%s) 需要更新， 证书Subject: %s, 证书名字：%s, 证书ID：%d", domain.Domain, cert.Subject, cert.Name, cert.CertID)
					err = setDCDNServerCertificateNotPanic(domain.Domain, cert.CertID, cert.Name)
					if err == nil {
						err := database.UpdateDCDNDomain(cert.CertID, cert.Name, cert.Subject, domain.Domain)
						if err != nil {
							logger.Errorf("在更新 域名 (%s) 状态到数据库时遇到了错误，但不影响后续检查: %s", domain.Domain, err.Error())
						}
					}
				} else {
					// 无需更新
					logger.Infof("确认域名 (%s) 无需更新证书，并找到证书，证书Subject: %s, 证书名字：%s, 证书ID：%d", domain.Domain, cert.Subject, cert.Name, cert.CertID)
				}
			} else {
				logger.Errorf("域名（%s）未知类型（%s）", domain.Domain, domain.Type)
			}
		}

		return nil
	} else if err != nil {
		return fmt.Errorf("aliyun cloud ssl cert/key upload error: %s", err.Error())
	} else {
		dbcerterr := database.UpdateCert(certID, certName, subject)
		if dbcerterr != nil {
			logger.Errorf("保存域名（%s）证书信息到数据库时发生了错误: %s", collection.Domain2Str(), dbcerterr.Error())
		}

		for _, domain := range collection.Domain {
			if domain.Type == config.DomainTypeCDN {
				err := setCDNServerCertificateNotPanic(domain.Domain, certID, certName)
				if dbcerterr != nil {
					logger.Errorf("因为证书信息未能写入数据库，所以 域名 (%s) 信息不尝试写入数据库", domain.Domain)
				} else if err == nil {
					err := database.UpdateCDNDomain(certID, certName, subject, domain.Domain)
					if err != nil {
						logger.Errorf("保存 域名 (%s) 信息到数据库发生了错误，但不影响后续域名更新: %s", domain.Domain, err.Error())
					}
				}
			} else if domain.Type == config.DomainTypeDCDN {
				err := setDCDNServerCertificateNotPanic(domain.Domain, certID, certName)
				if dbcerterr != nil {
					logger.Errorf("因为证书信息未能写入数据库，所以 域名 (%s) 信息不尝试写入数据库", domain.Domain)
				} else if err == nil {
					err := database.UpdateCDNDomain(certID, certName, subject, domain.Domain)
					if err != nil {
						logger.Errorf("保存 域名 (%s) 信息到数据库发生了错误，但不影响后续域名更新: %s", domain.Domain, err.Error())
					}
				}
			} else {
				logger.Errorf("域名（%s）未知类型（%s）", domain.Domain, domain.Type)
			}
		}
	}

	return nil
}

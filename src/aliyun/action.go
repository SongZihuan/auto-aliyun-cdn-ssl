package aliyun

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	cas "github.com/alibabacloud-go/cas-20200407/v3/client"
	cdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	dcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

func uploadCert(certData []byte, keyData []byte, resourceID string) (certID int64, name string, subject string, err error) {
	cert, err := utils.ReadCertificate(certData)
	if err != nil {
		return 0, "", "", fmt.Errorf("read cert error: %s", err.Error())
	}

	subject = utils.GetCertDomainSubject(cert)

	hash := sha256.Sum224(cert.RawSubjectPublicKeyInfo) // Sum256 太长
	fingerprint := hex.EncodeToString(hash[:])

	uploadUserCertificateRequest := &cas.UploadUserCertificateRequest{
		Name: tea.String(fingerprint),
		Cert: tea.String(string(certData)),
		Key:  tea.String(string(keyData)),
	}

	if resourceID != "" {
		uploadUserCertificateRequest.SetResourceGroupId(resourceID)
	}

	resp, tryErr := func() (resp *cas.UploadUserCertificateResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()

		return casClient.UploadUserCertificateWithOptions(uploadUserCertificateRequest, &util.RuntimeOptions{})
	}()
	if tryErr != nil {
		var sdkErr *tea.SDKError
		if errors.As(tryErr, &sdkErr) && tea.StringValue(sdkErr.Code) == "NameRepeat" {
			logger.Infof("证书已经存在, 证书名字：%s", fingerprint)
			return 0, fingerprint, subject, ErrCertExists
		}
		return 0, fingerprint, subject, tryErr
	}
	logger.Infof("上传成功, 证书Subject: %s, 证书名字：%s, 证书ID：%d, 请求ID：%s", subject, fingerprint, tea.Int64Value(resp.Body.CertId), tea.StringValue(resp.Body.RequestId))
	return tea.Int64Value(resp.Body.CertId), fingerprint, subject, nil
}

func setCDNServerCertificateNotPanic(domainName string, certID int64, certName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _err, ok := r.(error); ok {
				logger.Panicf("更新CDN 域名 (%s) 证书时发生了 不可预期的验证错误 被recover捕获，类型为error，错误消息是：%s", domainName, _err.Error())
				if err != nil {
					err = _err
				}
			} else {
				logger.Panicf("更新CDN 域名 (%s) 证书时发生了 不可预期的验证错误 被recover捕获，错误消息是：%v", domainName, r)
				if err != nil {
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	err = setCDNServerCertificate(domainName, certID, certName)
	if err != nil {
		logger.Infof("CDN加速 域名 (%s) 证书（%s）更新失败：%s", domainName, certName, err.Error())
		return err
	}

	return nil
}

func setCDNServerCertificate(domainName string, certID int64, certName string) (err error) {
	request := &cdn.SetCdnDomainSSLCertificateRequest{}
	request.DomainName = tea.String(strings.TrimPrefix(domainName, "*")) // 泛域名去除星号
	request.CertName = tea.String(certName)
	request.CertId = tea.Int64(certID)
	request.CertType = tea.String("cas")
	request.SSLProtocol = tea.String("on")
	if international {
		request.CertRegion = tea.String("ap-southeast-1") // 面向国际用农户
	} else {
		request.CertRegion = tea.String("cn-hangzhou") // 默认
	}

	_, tryErr := func() (resp *cdn.SetCdnDomainSSLCertificateResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()

		return cdnClient.SetCdnDomainSSLCertificate(request)
	}()
	if tryErr != nil {
		return tryErr
	}

	logger.Infof("CDN加速 域名（%s）证书（%s）更新成功, 并启用SSL", domainName, certName)
	return nil
}

func setDCDNServerCertificateNotPanic(domainName string, certID int64, certName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _err, ok := r.(error); ok {
				logger.Panicf("更新DCDN 域名 (%s) 证书时发生了 不可预期的验证错误 被recover捕获，类型为error，错误消息是：%s", domainName, _err.Error())
				if err != nil {
					err = _err
				}
			} else {
				logger.Panicf("更新DCDN 域名 (%s) 证书时发生了 不可预期的验证错误 被recover捕获，错误消息是：%v", domainName, r)
				if err != nil {
					err = fmt.Errorf("%v", r)
				}
			}
		}
	}()

	err = setDCDNServerCertificate(domainName, certID, certName)
	if err != nil {
		logger.Infof("DCDN加速 域名 (%s) 证书（%s）更新失败：%s", domainName, certName, err.Error())
		return err
	}

	return nil
}

func setDCDNServerCertificate(domainName string, certID int64, certName string) (err error) {
	request := &dcdn.SetDcdnDomainSSLCertificateRequest{}
	request.DomainName = tea.String(strings.TrimPrefix(domainName, "*")) // 泛域名去除星号
	request.CertName = tea.String(certName)
	request.CertId = tea.Int64(certID)
	request.CertType = tea.String("cas")
	request.SSLProtocol = tea.String("on")
	if international {
		request.CertRegion = tea.String("ap-southeast-1") // 面向国际用农户
	} else {
		request.CertRegion = tea.String("cn-hangzhou") // 默认
	}

	_, tryErr := func() (resp *dcdn.SetDcdnDomainSSLCertificateResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()

		return dcdnClient.SetDcdnDomainSSLCertificate(request)
	}()
	if tryErr != nil {
		return tryErr
	}

	logger.Infof("DCDN加速 域名（%s）证书（%s）更新成功, 并启用SSL", domainName, certName)
	return nil
}

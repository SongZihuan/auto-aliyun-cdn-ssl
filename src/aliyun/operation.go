package aliyun

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	cas "github.com/alibabacloud-go/cas-20200407/v3/client"
	cdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

func uploadCert(client *cas.Client, certData []byte, keyData []byte) (certID int64, name string, err error) {
	cert, err := utils.ReadCertificate(certData)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum224(cert.RawSubjectPublicKeyInfo) // Sum256 太长
	fingerprint := hex.EncodeToString(hash[:])

	uploadUserCertificateRequest := &cas.UploadUserCertificateRequest{
		Name: tea.String(fingerprint),
		Cert: tea.String(string(certData)),
		Key:  tea.String(string(keyData)),
	}

	resp, tryErr := func() (resp *cas.UploadUserCertificateResponse, err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()

		return client.UploadUserCertificateWithOptions(uploadUserCertificateRequest, &util.RuntimeOptions{})
	}()
	if tryErr != nil {
		var sdkErr *tea.SDKError
		if errors.As(tryErr, &sdkErr) && tea.StringValue(sdkErr.Code) == "NameRepeat" {
			logger.Warnf("证书已经存在, 证书名字：%s", fingerprint)
			return 0, fingerprint, ErrCertExists
		}
		return 0, fingerprint, tryErr
	}
	logger.Infof("上传成功, 证书名字：%s, 证书ID：%d, 请求ID：%s", fingerprint, tea.Int64Value(resp.Body.CertId), tea.StringValue(resp.Body.RequestId))
	return tea.Int64Value(resp.Body.CertId), fingerprint, nil
}

func setDomainServerCertificate(client *cdn.Client, domainName string, certID int64, certName string) (err error) {
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

		return client.SetCdnDomainSSLCertificate(request)
	}()
	if tryErr != nil {
		return tryErr
	}
	logger.Infof("CDN加速域名（%s）证书（%s）更新成功, 并启用SSL", domainName, certName)
	return nil
}

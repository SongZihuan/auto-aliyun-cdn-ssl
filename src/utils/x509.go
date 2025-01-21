package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

func ReadCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %s", err.Error())
	} else if cert == nil {
		return nil, fmt.Errorf("failed to parse certificate: return nil, unknown reason")
	}

	return cert, nil
}

func CheckCertWithDomain(cert *x509.Certificate, domain string) bool {
	// 遍历主题备用名称查找匹配的域名
	for _, name := range cert.DNSNames {
		if name == domain {
			return true // 找到了匹配的域名
		}
	}

	// 检查通用名作为回退，虽然现代实践倾向于使用SAN
	if cert.Subject.CommonName != "" && cert.Subject.CommonName == domain {
		return true // 通用名匹配
	}

	// 如果没有找到匹配，则返回错误
	return false
}

func CheckCertWithTime(cert *x509.Certificate, gracePeriod time.Duration) bool {
	now := time.Now()
	nowWithGracePeriod := now.Add(gracePeriod)

	if now.Before(cert.NotBefore) {
		return false
	} else if nowWithGracePeriod.After(cert.NotAfter) {
		return false
	}

	return true
}

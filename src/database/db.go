package database

import (
	"errors"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"gorm.io/gorm"
)

func UpdateCert(certID int64, name string, subject string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var cert CertRecord
		err := tx.Model(&CertRecord{}).Where("name = ?", name).Order("created_at desc").First(&cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cert = CertRecord{
				CertID:  certID,
				Name:    name,
				Subject: subject,
			}

			err = tx.Create(&cert).Error
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if cert.CertID != certID || cert.Name != name || cert.Subject != subject {
			cert.CertID = certID
			cert.Name = name
			cert.Subject = subject

			err = tx.Save(&cert).Error
			if err != nil {
				return err
			}

			err = tx.Model(&DomainRecord{}).Updates(map[string]interface{}{
				"cert_id": cert.CertID,
				"name":    cert.Name,
				"subject": cert.Subject,
			}).Error
			if err != nil {
				logger.Errorf("try update CDN domain SSL record because information does not match, but failed: %s", err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("create/update CAS cert record to SQLitte failed: %s", err.Error())
	}

	return nil
}

func UpdateDomain(certID int64, name string, subject string, domain string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var cert CertRecord
		err := tx.Model(&CertRecord{}).Where("name = ?", name).Order("created_at desc").First(&cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("update CDN domain SSL record to SQLite failed: cert record not found")
		} else if err != nil {
			return err
		} else if cert.CertID != certID || cert.Name != name || cert.Subject != subject {
			logger.Errorf("Update CDN domain SSL record to SQLite failed: information does not match (sqlite: cert-id=%d; name=%s; subject=%s) (be given: cert-id=%d; name=%s; subject=%s)", cert.CertID, cert.Name, cert.Subject, certID, name, subject)
		}

		record := DomainRecord{
			CertRecordID: cert.ID,
			CertID:       cert.CertID,
			Name:         cert.Name,
			Subject:      cert.Subject,
			Domain:       domain,
			CertUpdateAt: cert.UpdatedAt,
		}

		err = tx.Create(&record).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("create/update CDN domain SSL record to SQLite failed: %s", err.Error())
	}

	return nil
}

func CheckNeedUpdateDomain(certName string, domainName string) (cr *CertRecord, res bool, err error) {
	defer func() {
		if err != nil {
			res = false
			cr = nil
		}
	}()

	var cert CertRecord

	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&CertRecord{}).Where("name = ?", certName).Order("created_at desc").First(&cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check CDN domain SSL record from SQLite failed: cert record not found")
		} else if err != nil {
			return err
		}

		var domain DomainRecord
		err = tx.Model(&DomainRecord{}).Where("domain = ?", domainName).Order("created_at desc").First(&domain).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res = true
			return nil
		} else if err != nil {
			res = false
			return err
		} else if cert.UpdatedAt.After(domain.CertUpdateAt) {
			res = true
			return nil
		} else {
			res = false
			return nil
		}
	})
	if err != nil {
		return nil, false, fmt.Errorf("check CDN domain SSL record from SQLite failed: %s", err.Error())
	}

	return &cert, res, nil
}

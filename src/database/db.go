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

			err = tx.Model(&CDNDomainRecord{}).Updates(map[string]interface{}{
				"cert_id": cert.CertID,
				"name":    cert.Name,
				"subject": cert.Subject,
			}).Error
			if err != nil {
				logger.Errorf("try update CDN domain SSL record because information does not match, but failed: %s", err.Error())
			}

			err = tx.Model(&DCDNDomainRecord{}).Updates(map[string]interface{}{
				"cert_id": cert.CertID,
				"name":    cert.Name,
				"subject": cert.Subject,
			}).Error
			if err != nil {
				logger.Errorf("try update DCDN domain SSL record because information does not match, but failed: %s", err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("create/update CAS cert record to SQLitte failed: %s", err.Error())
	}

	return nil
}

func UpdateCDNDomain(certID int64, name string, subject string, domain string) error {
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

		record := CDNDomainRecord{
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

func CheckNeedUpdateCDNDomain(certName string, domainName string) (cert *CertRecord, need bool, err error) {
	defer func() {
		if err != nil || recover() != nil {
			need = false
			cert = nil
		}
	}()

	cert = new(CertRecord)

	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&CertRecord{}).Where("name = ?", certName).Order("created_at desc").First(cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cert = nil
			need = false
			return nil
		} else if err != nil {
			return err
		}

		var domain CDNDomainRecord
		err = tx.Model(&CDNDomainRecord{}).Where("domain = ?", domainName).Order("created_at desc").First(&domain).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			need = true
			return nil
		} else if err != nil {
			need = false
			return err
		} else if cert.UpdatedAt.After(domain.CertUpdateAt) {
			need = true
			return nil
		} else {
			need = false
			return nil
		}
	})
	if err != nil {
		return nil, false, fmt.Errorf("check CDN domain SSL record from SQLite failed: %s", err.Error())
	}

	return cert, need, nil
}

func UpdateDCDNDomain(certID int64, name string, subject string, domain string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var cert CertRecord
		err := tx.Model(&CertRecord{}).Where("name = ?", name).Order("created_at desc").First(&cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("update DCDN domain SSL record to SQLite failed: cert record not found")
		} else if err != nil {
			return err
		} else if cert.CertID != certID || cert.Name != name || cert.Subject != subject {
			logger.Errorf("Update DCDN domain SSL record to SQLite failed: information does not match (sqlite: cert-id=%d; name=%s; subject=%s) (be given: cert-id=%d; name=%s; subject=%s)", cert.CertID, cert.Name, cert.Subject, certID, name, subject)
		}

		record := DCDNDomainRecord{
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
		return fmt.Errorf("create/update DCDN domain SSL record to SQLite failed: %s", err.Error())
	}

	return nil
}

func CheckNeedUpdateDCDNDomain(certName string, domainName string) (cert *CertRecord, need bool, err error) {
	defer func() {
		if err != nil || recover() != nil {
			need = false
			cert = nil
		}
	}()

	cert = new(CertRecord)

	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&CertRecord{}).Where("name = ?", certName).Order("created_at desc").First(cert).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cert = nil
			need = false
			return nil
		} else if err != nil {
			return err
		}

		var domain DCDNDomainRecord
		err = tx.Model(&DCDNDomainRecord{}).Where("domain = ?", domainName).Order("created_at desc").First(&domain).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			need = true
			return nil
		} else if err != nil {
			need = false
			return err
		} else if cert.UpdatedAt.After(domain.CertUpdateAt) {
			need = true
			return nil
		} else {
			need = false
			return nil
		}
	})
	if err != nil {
		return nil, false, fmt.Errorf("check DCDN domain SSL record from SQLite failed: %s", err.Error())
	}

	return cert, need, nil
}

package database

import "gorm.io/gorm"

type CertRecord struct {
	gorm.Model
	CertID  int64  `gorm:"column:cert_id;not null;uniqueIndex:unq_idx_cert"`
	Name    string `gorm:"column:name;type:VARCHAR(100);not null;uniqueIndex:unq_idx_cert"`
	Subject string `gorm:"column:subject;type:VARCHAR(100);not null;"`
}

func (*CertRecord) TableName() string {
	return "cert_record"
}

type DomainRecord struct {
	gorm.Model
	CertRecordID uint   `gorm:"column:cert_record_id;not null;"`
	CertID       int64  `gorm:"column:cert_id;not null"`
	Name         string `gorm:"column:name;type:VARCHAR(100);not null"`
	Subject      string `gorm:"column:subject;type:VARCHAR(100);not null;"`
	Domain       string `gorm:"type:VARCHAR(100);not null;"` // 允许多次重复记录
}

func (*DomainRecord) TableName() string {
	return "domain_record"
}

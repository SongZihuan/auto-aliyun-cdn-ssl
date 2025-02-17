package database

import (
	"gorm.io/gorm"
	"time"
)

// Model gorm.Model的仿写，明确了键名
type Model struct {
	ID        uint           `gorm:"column:id;primarykey"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

type CertRecord struct {
	Model
	CertID  int64  `gorm:"column:cert_id;not null;uniqueIndex:unq_idx_cert"`
	Name    string `gorm:"column:name;type:VARCHAR(100);not null;uniqueIndex:unq_idx_cert"`
	Subject string `gorm:"column:subject;type:VARCHAR(100);not null;"`
}

func (*CertRecord) TableName() string {
	return "cert_record"
}

type CDNDomainRecord struct {
	Model
	CertRecordID uint      `gorm:"column:cert_record_id;not null;"`
	CertID       int64     `gorm:"column:cert_id;not null"`
	Name         string    `gorm:"column:name;type:VARCHAR(100);not null"`
	Subject      string    `gorm:"column:subject;type:VARCHAR(100);not null;"`
	Domain       string    `gorm:"type:VARCHAR(100);not null;"` // 允许多次重复记录
	CertUpdateAt time.Time `gorm:"column:cert_update_time;not null"`
}

func (*CDNDomainRecord) TableName() string {
	return "cdn_domain_record"
}

type DCDNDomainRecord struct {
	Model
	CertRecordID uint      `gorm:"column:cert_record_id;not null;"`
	CertID       int64     `gorm:"column:cert_id;not null"`
	Name         string    `gorm:"column:name;type:VARCHAR(100);not null"`
	Subject      string    `gorm:"column:subject;type:VARCHAR(100);not null;"`
	Domain       string    `gorm:"type:VARCHAR(100);not null;"` // 允许多次重复记录
	CertUpdateAt time.Time `gorm:"column:cert_update_time;not null"`
}

func (*DCDNDomainRecord) TableName() string {
	return "dcdn_domain_record"
}

package database

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitSQLite() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	sqlfilepath := config.GetConfig().SQLFilePath

	_db, err := gorm.Open(sqlite.Open(sqlfilepath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to sqlite (%s) failed: %s", sqlfilepath, err)
	}

	err = _db.AutoMigrate(&CertRecord{}, &CDNDomainRecord{}, &DCDNDomainRecord{})
	if err != nil {
		return fmt.Errorf("migrate sqlite (%s) failed: %s", sqlfilepath, err)
	}

	db = _db
	return nil
}

func CloseSQLite() {
	if db == nil {
		return
	}

	defer func() {
		db = nil
	}()

	if !config.IsReady() {
		panic("config is not ready")
	}

	if config.GetConfig().ActiveShutdown.IsEnable(false) {
		// https://github.com/go-gorm/gorm/issues/3145
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

package server

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/aliyun"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
)

func Server() error {
	cfg := config.GetConfig().DomainConfig

	logger.Info("Server start...")
	for _, domain := range cfg.Domains {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						logger.Panicf("aliyun update CDN HTTOS by domain (%s) panic: %s", domain, err.Error())
					} else {
						logger.Panicf("aliyun update CDN HTTOS by domain (%s) panic: %v", domain, r)
					}
				}
			}()

			certPath, prikeyPath := domain.GetFilePath()
			err := aliyun.UpdateCDNHttpsByFilePath(domain.Domain, certPath, prikeyPath)
			if err != nil {
				logger.Errorf("aliyun update CDN HTTOS by domain (%s) error: %s", domain, err.Error())
			}
		}()
	}
	logger.Info("Server finish...")

	return nil
}

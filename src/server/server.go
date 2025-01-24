package server

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/aliyun"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"strings"
)

func Server() error {
	cfg := config.GetConfig().DomainListsGroup

	logger.Info("Server start...")
	for index, collection := range cfg.Collection {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						logger.Panicf("aliyun update CDN HTTPS by domains/collection (%s / %d) panic: %s", strings.Join(collection.Domain, ", "), index, err.Error())
					} else {
						logger.Panicf("aliyun update CDN HTTPS by domains/collection (%s / %d) panic: %v", strings.Join(collection.Domain, ", "), index, r)
					}
				}
			}()

			certPath, prikeyPath := collection.GetFilePath()
			err := aliyun.UpdateCDNHttpsByFilePath(collection.Domain, certPath, prikeyPath)
			if err != nil {
				logger.Panicf("aliyun update CDN HTTPS by domains/collection (%s / %d) panic: %s", strings.Join(collection.Domain, ", "), index, err.Error())
			}
		}()
	}
	logger.Info("Server finish...")

	return nil
}

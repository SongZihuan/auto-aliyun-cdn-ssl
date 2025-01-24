package server

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/aliyun"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"strings"
)

func Server() error {
	cfg := config.GetConfig().DomainListsGroup

	logger.Info("服务开始...")
	for index, collection := range cfg.Collection {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						logger.Panicf("本集合的域名在进行更新时遇到了 不可预期的严重错误 被recover捕获，类型为error，其他集合继续运行 (域名：%s / 序号：%d) 错误消息: %s", strings.Join(collection.Domain, ", "), index, err.Error())
					} else {
						logger.Panicf("本集合的域名在进行更新时遇到了 不可预期的严重错误 被recover捕获，其他集合继续运行 (域名：%s / 序号：%d) 错误消息: %v", strings.Join(collection.Domain, ", "), index, r)
					}
				}
			}()

			logger.Infof("开始第 %d 组集合（从0算起）的更新服务，包括域名：%s", index, strings.Join(collection.Domain, ", "))
			certPath, prikeyPath := collection.GetFilePath()
			logger.Infof("获取到证书路径：%s，密钥路径：%s", certPath, prikeyPath)
			err := aliyun.UpdateCDNHttpsByFilePath(collection.Domain, certPath, prikeyPath)
			if err != nil {
				logger.Errorf("本集合的域名在进行更新时遇到了 错误 被recover捕获，类型为error，其他集合继续运行 (域名：%s / 序号：%d) 错误消息: %s", strings.Join(collection.Domain, ", "), index, err.Error())
			}
			logger.Infof("完成第 %d 组集合（从0算起）的更新服务，包括域名：%s", index, strings.Join(collection.Domain, ", "))
		}()
	}
	logger.Info("服务结束...")

	return nil
}

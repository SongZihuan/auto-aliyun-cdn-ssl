package config

import (
	"os"
)

func InitConfig(configPath string) ConfigError {
	var err error
	config, err = newConfig(configPath)
	if err != nil {
		return NewConfigError(err.Error())
	}

	cfgErr := config.Init()
	if cfgErr != nil && cfgErr.IsError() {
		return cfgErr
	}

	if !config.IsReady() {
		return NewConfigError("config not ready")
	}

	return nil
}

func IsReady() bool {
	return config.IsReady()
}

func GetConfig() *YamlConfig {
	return config.GetConfig()
}

func GetSignalChan() chan os.Signal {
	return config.GetSignalChan()
}

var config *ConfigStruct

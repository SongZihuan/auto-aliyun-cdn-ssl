package config

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"os"
)

const EnvPrefix = "ALIYUN_"
const (
	EnvAliyunKey    = EnvPrefix + "KEY"
	EnvAliyunSecret = EnvPrefix + "SECRET"
)

type AliyunConfig struct {
	Key           string           `yaml:"key"`
	Secret        string           `yaml:"secret"`
	International utils.StringBool `yaml:"international"`
	ResourceID    string           `yaml:"resource-id"`
}

func (a *AliyunConfig) SetDefault() {
	if a.Key == "" {
		a.Key = os.Getenv(EnvAliyunKey)
	}

	if a.Key == "" {
		a.Key = os.Getenv(EnvAliyunSecret)
	}

	a.International.SetDefaultDisable()
	return
}

func (a *AliyunConfig) Check() ConfigError {
	if a.Key == "" || a.Secret == "" {
		return NewConfigError("aliyun key or secret is empty")
	}

	return nil
}

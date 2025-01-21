package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type YamlConfig struct {
	GlobalConfig `yaml:",inline"`
	DomainConfig `yaml:",inline"`

	Aliyun AliyunConfig `yaml:"aliyun"`
}

func (y *YamlConfig) Init() error {
	return nil
}

func (y *YamlConfig) SetDefault(configPath string) {
	y.GlobalConfig.SetDefault()
	y.DomainConfig.SetDefault(configPath)
	y.Aliyun.SetDefault()
}

func (y *YamlConfig) Check() (err ConfigError) {
	err = y.GlobalConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.DomainConfig.Check()
	if err != nil && err.IsError() {
		return err
	}

	err = y.Aliyun.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (y *YamlConfig) Parser(filepath string) ParserError {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return NewParserError(err, err.Error())
	}

	err = yaml.Unmarshal(file, y)
	if err != nil {
		return NewParserError(err, err.Error())
	}

	return nil
}

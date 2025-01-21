package config

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/flagparser"
	"os"
	"sync"
)

type ConfigStruct struct {
	ConfigLock sync.Mutex

	configReady   bool
	yamlHasParser bool
	sigchan       chan os.Signal
	configPath    string

	Yaml *YamlConfig
}

func newConfig(configPath string) (*ConfigStruct, error) {
	if configPath == "" {
		if !flagparser.IsReady() {
			panic("flag is not ready")
		}
		configPath = flagparser.ConfigFile()
	}

	return &ConfigStruct{
		configPath:    configPath,
		configReady:   false,
		yamlHasParser: false,
		sigchan:       make(chan os.Signal),
		Yaml:          nil,
	}, nil
}

func (c *ConfigStruct) Init() (err ConfigError) {
	if c.IsReady() { // 使用IsReady而不是isReady，确保上锁
		return
	}

	initErr := c.init()
	if initErr != nil {
		return NewConfigError("init error: " + initErr.Error())
	}

	parserErr := c.Parser(c.configPath)
	if parserErr != nil {
		return NewConfigError("parser error: " + parserErr.Error())
	} else if !c.yamlHasParser {
		return NewConfigError("parser error: unknown")
	}

	c.SetDefault()

	err = c.Check()
	if err != nil && err.IsError() {
		return err
	}

	c.configReady = true
	return nil
}

func (c *ConfigStruct) Parser(filepath string) ParserError {
	err := c.Yaml.Parser(filepath)
	if err != nil {
		return err
	}

	c.yamlHasParser = true
	return nil
}

func (c *ConfigStruct) SetDefault() {
	if !c.yamlHasParser {
		panic("yaml must parser first")
	}

	c.Yaml.SetDefault(c.configPath)
}

func (c *ConfigStruct) Check() (err ConfigError) {
	err = c.Yaml.Check()
	if err != nil && err.IsError() {
		return err
	}

	return nil
}

func (c *ConfigStruct) isReady() bool {
	return c.yamlHasParser && c.configReady
}

func (c *ConfigStruct) init() error {
	c.configReady = false
	c.yamlHasParser = false

	err := initSignal(c.sigchan)
	if err != nil {
		return err
	}

	c.Yaml = new(YamlConfig)
	err = c.Yaml.Init()
	if err != nil {
		return err
	}

	return nil
}

// export func

func (c *ConfigStruct) IsReady() bool {
	return c.isReady()
}

func (c *ConfigStruct) GetSignalChan() chan os.Signal {
	if !c.isReady() {
		panic("config is not ready")
	}

	return c.sigchan
}

func (c *ConfigStruct) GetConfig() *YamlConfig {
	if !c.isReady() {
		panic("config is not ready")
	}

	return c.Yaml
}

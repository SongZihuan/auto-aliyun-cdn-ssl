package config

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"os"
)

const EnvModeName = "HUAN_PROXY_MODE"

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

type LoggerLevel string

var levelMap = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
	"panic": true,
	"none":  true,
}

type GlobalConfig struct {
	Mode     string           `yaml:"mode"`
	LogLevel string           `yaml:"loglevel"`
	LogTag   utils.StringBool `yaml:"logtag"`
}

func (g *GlobalConfig) SetDefault() {
	if g.Mode == "" {
		g.Mode = os.Getenv(EnvModeName)
	}

	if g.Mode == "" {
		g.Mode = DebugMode
	}

	_ = os.Setenv(EnvModeName, g.Mode)

	if g.LogLevel == "" && (g.Mode == DebugMode || g.Mode == TestMode) {
		g.LogLevel = "debug"
	} else if g.LogLevel == "" {
		g.LogLevel = "warn"
	}

	if g.Mode == DebugMode || g.Mode == TestMode {
		g.LogTag.SetDefaultEnable()
	} else {
		g.LogTag.SetDefaultDisable()
	}

	return
}

func (g *GlobalConfig) Check() ConfigError {
	if g.Mode != DebugMode && g.Mode != ReleaseMode && g.Mode != TestMode {
		return NewConfigError("bad mode")
	}

	if _, ok := levelMap[g.LogLevel]; !ok {
		return NewConfigError("log level error")
	}

	return nil
}

func (g *GlobalConfig) GetRunMode() string {
	return g.Mode
}

func (g *GlobalConfig) IsDebug() bool {
	return g.Mode == DebugMode
}

func (g *GlobalConfig) IsRelease() bool {
	return g.Mode == ReleaseMode
}

func (g *GlobalConfig) IsTest() bool {
	return g.Mode == TestMode
}

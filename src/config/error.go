package config

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
)

type ConfigError interface {
	error
	Msg() string
	Error() string
	Warning() string
	IsError() bool
	IsWarning() bool
}

func NewConfigError(msg string) ConfigError {
	fmt.Println(utils.FormatTextToWidth(fmt.Sprintf("config error: %s", msg), utils.NormalConsoleWidth))
	return &configError{msg: msg, isError: true}
}

func NewConfigWarning(msg string) ConfigError {
	fmt.Println(utils.FormatTextToWidth(fmt.Sprintf("config warning: %s", msg), utils.NormalConsoleWidth))
	return &configError{msg: msg, isError: false}
}

type configError struct {
	msg     string
	isError bool
}

func (e *configError) Msg() string {
	if e.isError {
		return "config error: " + e.Error()
	}
	return "config warning: " + e.Warning()
}

func (e *configError) Error() string {
	return e.msg
}

func (e *configError) Warning() string {
	return e.msg
}

func (e *configError) IsError() bool {
	return e.isError
}

func (e *configError) IsWarning() bool {
	return !e.isError
}

type ParserError interface {
	error
	Error() string
	Data() interface{}
}

type parserError struct {
	msg  string
	data interface{}
}

func NewParserError(data interface{}, msg ...string) ParserError {
	if len(msg) == 1 {
		return &parserError{msg[0], data}
	}
	return &parserError{"config parser error: " + fmt.Sprint(data), data}
}

func WarpParserError(err error) ParserError {
	return &parserError{"config parser error: " + err.Error(), err}
}

func (e *parserError) Error() string {
	return e.msg
}

func (e *parserError) Data() interface{} {
	return e.data
}

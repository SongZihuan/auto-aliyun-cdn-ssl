package flagparser

import "fmt"

type FlagError interface {
	error
	Error() string
	Data() interface{}
}

type flagError struct {
	msg  string
	data interface{}
}

func NewFlagError(data interface{}, msg ...string) FlagError {
	if len(msg) == 1 {
		return &flagError{msg[0], data}
	}
	return &flagError{"flag error: " + fmt.Sprint(data), data}
}

func (e *flagError) Error() string {
	return e.msg
}

func (e *flagError) Data() interface{} {
	return e.data
}

package utils

import (
	"os"
)

func ExitByError(err error, code ...int) int {
	if err == nil {
		return ExitByErrorMsg("")
	} else {
		return ExitByErrorMsg(err.Error(), code...)
	}
}

func ExitByErrorMsg(msg string, code ...int) int {
	if len(msg) == 0 {
		msg = "exit: unknown error"
	}

	return ErrorExit(msg, code...)
}

func ErrorExit(msg string, code ...int) int {
	if len(msg) == 0 {
		SayGoodByef("%s", "Encountered an error, abnormal offline/shutdown.")
	} else {
		SayGoodByef("Encountered an error, abnormal offline/shutdown: %s\n", msg)
	}

	if len(code) == 1 && code[0] != 0 {
		return Exit(code[0])
	} else {
		return Exit(1)
	}
}

func Exit(code ...int) int {
	if len(code) == 1 {
		os.Exit(code[0])
		return code[0]
	} else {
		os.Exit(0)
		return 0
	}
}

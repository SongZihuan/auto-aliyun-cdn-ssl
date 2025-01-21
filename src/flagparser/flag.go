package flagparser

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"os"
)

var isReady = false

func IsReady() bool {
	return data.isReady() && isReady
}

var StopFlag = fmt.Errorf("stop")

func InitFlag() (err error) {
	if isReady {
		return nil
	}

	defer func() {
		if e := recover(); e != nil {
			err = NewFlagError(e)
			return
		}
	}()

	initData()

	SetOutput(os.Stdout)

	var hasPrint = false

	if Version() {
		_, _ = PrintVersion()
		hasPrint = true
	}

	if License() {
		if hasPrint {
			_, _ = PrintLF()
		}
		_, _ = PrintLicense()
		hasPrint = true
	}

	if Report() {
		if hasPrint {
			_, _ = PrintLF()
		}
		_, _ = PrintReport()
	}

	if Help() {
		if hasPrint {
			_, _ = PrintLF()
		}
		_, _ = PrintUsage()
		hasPrint = true
	}

	if NotRunMode() {
		return StopFlag
	}

	err = checkFlag()
	if err != nil {
		return err
	}

	isReady = true
	return nil
}

func checkFlag() error {
	if !utils.IsExists(ConfigFile()) {
		return fmt.Errorf("config file not exists")
	}

	return nil
}

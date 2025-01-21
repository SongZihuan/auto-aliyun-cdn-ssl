package mainfunc

import (
	"errors"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/aliyun"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/config"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/flagparser"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/logger"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/server"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"os"
)

func MainV1() int {
	var err error

	err = flagparser.InitFlag()
	if errors.Is(err, flagparser.StopFlag) {
		return 0
	} else if err != nil {
		return utils.ExitByError(err)
	}

	if !flagparser.IsReady() {
		return utils.ExitByErrorMsg("flag parser unknown error")
	}

	utils.SayHellof("%s", "The backend service program starts normally, thank you.")
	defer func() {
		utils.SayGoodByef("%s", "The backend service program is offline/shutdown normally, thank you.")
	}()

	cfgErr := config.InitConfig(flagparser.ConfigFile())
	if cfgErr != nil && cfgErr.IsError() {
		return utils.ExitByError(cfgErr)
	}

	if !config.IsReady() {
		return utils.ExitByErrorMsg("config parser unknown error")
	}

	err = logger.InitLogger(os.Stdout, os.Stderr)
	if err != nil {
		return utils.ExitByError(err)
	}

	if !logger.IsReady() {
		return utils.ExitByErrorMsg("logger unknown error")
	}

	logger.Executablef("%s", "ready")
	logger.Infof("run mode: %s", config.GetConfig().GlobalConfig.GetRunMode())

	err = aliyun.Init()
	if err != nil {
		return utils.ExitByError(err)
	}

	err = server.Server()
	if err != nil {
		return utils.ExitByError(err)
	}

	return 0
}

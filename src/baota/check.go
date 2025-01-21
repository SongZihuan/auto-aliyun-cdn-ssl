package baota

import (
	"bytes"
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// 宝塔 Let's Encrypt 证书保存位置：/www/server/panel/vhost/letsencrypt
const BT_LETSENCRYPT_SSL = "/www/server/panel/vhost/letsencrypt"

var isBaota *bool

func IsLinuxBaoTa() (res bool) {
	if isBaota != nil {
		return *isBaota
	}

	defer func() {
		isBaota = &res
	}()

	if runtime.GOOS != "linux" {
		fmt.Println("运行环境：不是 宝塔 面板")
		return false
	}

	var stdout, stderr bytes.Buffer

	cmd := exec.Command("pgrep", "-f", "BT-Panel")
	cmd.Stdout = &stdout // 将 stdout 重定
	cmd.Stdout = &stderr // 将 stderr 重定
	err := cmd.Run()
	if err != nil {
		fmt.Println("运行环境：不是 宝塔 面板")
		return false
	}

	_stdout := utils.StringToOnlyPrint(stdout.String())
	_stderr := utils.StringToOnlyPrint(stderr.String())

	if _stderr == "" || _stdout != "" {
		fmt.Println("运行环境：Linux 宝塔版本")
		return true
	}

	if utils.IsDir("/www/server/panel/BTPanel") {
		fmt.Println("运行环境：Linux 宝塔版本")
		return true
	}

	initsh, err := os.ReadFile("/www/server/panel/init.sh")
	if err != nil {
		fmt.Println("运行环境：不是 宝塔 面板")
		return false
	}

	if strings.Contains(string(initsh), "宝塔") {
		fmt.Println("运行环境：Linux 宝塔版本")
		return true
	}

	licensetxt, err := os.ReadFile("/www/server/panel/license.txt")
	if err != nil {
		fmt.Println("运行环境：不是 宝塔 面板")
		return false
	}

	if strings.Contains(string(licensetxt), "宝塔") {
		fmt.Println("运行环境：Linux 宝塔版本")
		return true
	}

	fmt.Println("运行环境：不是 宝塔 面板")
	return false
}

func GetBaoTaLetsEncryptDir() string {
	return BT_LETSENCRYPT_SSL
}

func HasBaoTaLetsEncrypt() bool {
	if !IsLinuxBaoTa() {
		return false
	}

	return utils.IsDir(GetBaoTaLetsEncryptDir())
}

package config

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/baota"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"path"
	"path/filepath"
)

type DomainListCollection struct {
	Domain []string `yaml:"domain"`
	Dir    string   `yaml:"dir"`
	Cert   string   `yaml:"cert"`
	Key    string   `yaml:"pri"`
}

type DomainListsGroup struct {
	SQLFilePath    string                 `yaml:"sqlfilepath"`
	ActiveShutdown utils.StringBool       `yaml:"activeshutdown"`
	RootDir        string                 `yaml:"rootrir"`
	Collection     []DomainListCollection `yaml:"collection"`
}

func (d *DomainListsGroup) SetDefault(configPath string) {
	if d.RootDir == "" {
		if baota.HasBaoTaLetsEncrypt() {
			d.RootDir = baota.GetBaoTaLetsEncryptDir()
		} else {
			d.RootDir = path.Dir(configPath)
		}
	}

	d.ActiveShutdown.SetDefaultDisable()

	if d.SQLFilePath == "" {
		if baota.HasBaoTaLetsEncrypt() {
			d.SQLFilePath = path.Join(configPath, "auto-aliyun-cdn-ssl.db")
		} else {
			d.SQLFilePath = "./auto-aliyun-cdn-ssl.db"
		}
	}

	return
}

func (d *DomainListsGroup) Check() ConfigError {
	if !utils.IsDir(d.RootDir) {
		return NewConfigError("root dir is not a dir")
	}

	for _, domainLst := range d.Collection {
		if len(domainLst.Domain) == 0 {
			return NewConfigError("domain list is empty")
		}

		for _, domain := range domainLst.Domain {
			if domain == "" {
				return NewConfigError("domain is empty")
			} else if !utils.IsValidDomain(domain) && !utils.IsValidWildcardDomain(domain) {
				return NewConfigError("domain is not valid")
			}
		}
	}

	return nil
}

func (d *DomainListCollection) GetFilePath() (certPath string, prikeyPath string) {
	if !IsReady() {
		panic("config is not ready")
	}

	rootDir := GetConfig().RootDir

	if baota.IsLinuxBaoTa() && rootDir != baota.GetBaoTaLetsEncryptDir() {
		fmt.Printf("提示：运行在宝塔环境，但非Let's Encrypt目录")
	}

	if d.Dir != "" {
		if filepath.IsAbs(d.Dir) {
			rootDir = d.Dir
		} else {
			rootDir = path.Join(rootDir, d.Dir)
		}
	} else {
		rootDir = path.Join(rootDir, d.Domain[0])
	}

	if rootDir == "" {
		panic("root dir is empty") // 实际上这不可能发生，因为d.RootDir是非空（SetDefault），而d.Dir只有非空是才会被赋值
	}

	certPath = path.Join(rootDir, "fullchain.pem")
	if d.Cert != "" {
		certPath = path.Join(rootDir, d.Cert)
	}

	prikeyPath = path.Join(rootDir, "privkey.pem")
	if d.Key != "" {
		prikeyPath = path.Join(rootDir, d.Key)
	}

	return certPath, prikeyPath
}

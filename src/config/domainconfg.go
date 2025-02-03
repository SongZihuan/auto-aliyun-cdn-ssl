package config

import (
	"fmt"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/baota"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"path"
	"path/filepath"
	"strings"
)

const (
	DomainTypeCDN  = "cdn"
	DomainTypeDCDN = "dcdn"
)

type Domain struct {
	Domain string `yaml:"domain"`
	Type   string `yaml:"type"`
}

type DomainListCollection struct {
	Domain []*Domain `yaml:"domain"`
	Dir    string    `yaml:"dir"`
	Cert   string    `yaml:"cert"`
	Key    string    `yaml:"pri"`
}

type DomainListsGroup struct {
	SQLFilePath    string                  `yaml:"sqlfilepath"`
	ActiveShutdown utils.StringBool        `yaml:"activeshutdown"`
	RootDir        string                  `yaml:"rootrir"`
	Collection     []*DomainListCollection `yaml:"collection"`
}

func (c *DomainListCollection) Domain2Str() string {
	var builder strings.Builder
	var seq = ", "

	for _, domain := range c.Domain {
		builder.WriteString(fmt.Sprintf("%s（%s）%s", domain.Domain, domain.Type, seq))
	}

	res := strings.TrimRight(builder.String(), seq)
	return res
}

func (d *DomainListsGroup) SetDefault(configPath string) {
	if d.SQLFilePath == "" {
		if baota.HasBaoTaLetsEncrypt() {
			d.SQLFilePath = path.Join(configPath, "auto-aliyun-cdn-ssl.db")
		} else {
			d.SQLFilePath = "./auto-aliyun-cdn-ssl.db"
		}
	}

	d.ActiveShutdown.SetDefaultDisable()

	if d.RootDir == "" {
		if baota.HasBaoTaLetsEncrypt() {
			d.RootDir = baota.GetBaoTaLetsEncryptDir()
		} else {
			d.RootDir = path.Dir(configPath)
		}
	}

	for _, c := range d.Collection {
		for _, domain := range c.Domain {
			domain.Domain = strings.TrimSpace(strings.ToLower(domain.Domain))
			domain.Type = strings.TrimSpace(strings.ToLower(domain.Type))
			if domain.Type == "" {
				domain.Type = DomainTypeCDN
			}
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
			if domain.Domain == "" {
				return NewConfigError("domain is empty")
			} else if !utils.IsValidDomain(domain.Domain) && !utils.IsValidWildcardDomain(domain.Domain) {
				return NewConfigError("domain is not valid")
			}

			if domain.Type != DomainTypeCDN && domain.Type != DomainTypeDCDN {
				return NewConfigError("domain type is not valid")
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
		rootDir = path.Join(rootDir, d.Domain[0].Domain)
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

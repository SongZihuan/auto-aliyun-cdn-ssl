package config

import (
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/baota"
	"github.com/SongZihuan/auto-aliyun-cdn-ssl/src/utils"
	"path"
)

type Domain struct {
	Domain  string `yaml:"domain"`
	RootDir string `yaml:"-"`
	Dir     string `yaml:"dir"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"pri"`
}

type DomainConfig struct {
	RootDir string   `yaml:"rootrir"`
	Domains []Domain `yaml:"domains"`
}

func (d *DomainConfig) SetDefault(configPath string) {
	if d.RootDir == "" {
		if baota.HasBaoTaLetsEncrypt() {
			d.RootDir = baota.GetBaoTaLetsEncryptDir()
		} else {
			d.RootDir = path.Dir(configPath)
		}
	}

	for _, domain := range d.Domains {
		domain.RootDir = d.RootDir
	}

	return
}

func (d *DomainConfig) Check() ConfigError {
	if !utils.IsDir(d.RootDir) {
		return NewConfigError("root dir is not a dir")
	}

	for _, domain := range d.Domains {
		if domain.Domain == "" {
			return NewConfigError("domain is empty")
		} else if !utils.IsValidDomain(domain.Domain) && !utils.IsValidWildcardDomain(domain.Domain) {
			return NewConfigError("domain is not valid")
		}
	}

	return nil
}

func (d *Domain) GetFilePath() (certPath string, prikeyPath string) {
	rootDir := d.RootDir

	if d.Dir != "" {
		if path.IsAbs(d.Dir) {
			rootDir = d.Dir
		} else {
			rootDir = path.Join(rootDir, d.Dir)
		}
	} else {
		rootDir = path.Join(rootDir, d.Domain)
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

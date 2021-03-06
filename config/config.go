package config

import (
	"github.com/csumissu/SkyDisk/util"
	"gopkg.in/ini.v1"
)

type database struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Name     string `ini:"dbname"`
}

type server struct {
	Port                 int   `ini:"port"`
	MultipartMaxMemoryMB int64 `ini:"multipartMaxMemoryMB"`
	MultipartMaxUploadMB int64 `ini:"multipartMaxUploadMB"`
}

type cors struct {
	AllowOrigins     []string `ini:"allowOrigins"`
	AllowMethods     []string `ini:"allowMethods"`
	AllowHeaders     []string `ini:"allowHeaders"`
	AllowCredentials bool     `ini:"allowCredentials"`
	ExposeHeaders    []string `ini:"exposeHeaders"`
}

type redis struct {
	Network  string `ini:"network"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	DB       int    `ini:"db"`
}

type jwt struct {
	SigningKey      string `ini:"signingKey"`
	ExpirationHours int    `ini:"expirationHours"`
}

func init() {
	cfg, err := ini.Load("config/resources/dev.ini")
	if err != nil {
		util.Logger.Panic("fail to read ini file: %v", err)
	}

	loadSections(cfg)
}

func loadSections(cfg *ini.File) {
	sections := map[string]interface{}{
		"database": DatabaseCfg,
		"server":   ServerCfg,
		"cors":     CORSCfg,
		"redis":    RedisCfg,
		"jwt":      JwtCfg,
	}

	for sectionName, sectionStruct := range sections {
		err := cfg.Section(sectionName).MapTo(sectionStruct)
		if err != nil {
			util.Logger.Panic("fail to extract %s section from ini file, %v", sectionName, err)
		}
	}
}

package conf

import (
	"fmt"

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
	Port int `ini:"port"`
}

type cors struct {
	AllowOrigins     []string `ini:"allowOrigins"`
	AllowMethods     []string `ini:"allowMethods"`
	AllowHeaders     []string `ini:"allowHeaders"`
	AllowCredentials bool     `ini:"allowCredentials"`
	ExposeHeaders    []string `ini:"exposeHeaders"`
}

func init() {
	cfg, err := ini.Load("conf/dev.ini")
	if err != nil {
		panic(fmt.Sprintf("Fail to read ini file: %v", err))
	}

	loadSections(cfg)
}

func loadSections(cfg *ini.File) {
	sections := map[string]interface{}{
		"database": DatabaseCfg,
		"server":   ServerCfg,
		"cors":     CORSCfg,
	}

	for sectionName, sectionStruct := range sections {
		err := cfg.Section(sectionName).MapTo(sectionStruct)
		if err != nil {
			panic(fmt.Sprintf("Fail to extract %s section from ini file, %v", sectionName, err))
		}
	}
}

package config

import (
	"github.com/go-ini/ini"
	"log"
)

type Database struct {
	Name string
}

type Server struct {
	Port   uint
	Domain string
}

var DatabaseSetting = &Database{}
var ServerSetting = &Server{}

var cfg *ini.File

func init() {
	var err error
	//str, _ := os.Getwd()
	//println(str)
	cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("server", ServerSetting)
}

func mapTo(name string, v interface{}) {
	err := cfg.Section(name).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", name, err)
	}
}

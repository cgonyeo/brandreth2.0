package config

import (
	"code.google.com/p/gcfg"
	golog "github.com/op/go-logging"
)

var log = golog.MustGetLogger("main")

var Config = struct {
	Templates struct {
		Path string
	}
	Admins struct {
		Name []string
	}
}{}

func LoadConfig(filename string) {
	err := gcfg.ReadFileInto(&Config, filename)
	if err != nil {
		log.Error("Error parsing config file: %v", err)
	}
}

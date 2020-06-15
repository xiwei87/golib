package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg Config

type Config struct {
	Server ServerCfg `yaml:"server"`
	Http   HttpCfg   `yaml:"http"`
	Sqlite SqliteCfg `yaml:"sqlite"`
}

type ServerCfg struct {
	LogPath  string `yaml:"LogPath"`
	LogLevel string `yaml:"LogLevel"`
	LogSave  int    `yaml:"LogSave"`
	MaxCpus  int    `yaml:"MaxCpus"`
}

type HttpCfg struct {
	ListenPort       int  `yaml:"ListenPort"`
	ReadTimeout      uint `yaml:"ReadTimeout"`
	ReadIdleTimeout  uint `yaml:"ReadIdleTimeout"`
	WriteTimeout     uint `yaml:"WriteTimeout"`
	WriteIdleTimeout uint `yaml:"WriteIdleTimeout"`
	MaxHeaderSize    int  `yaml:"MaxHeaderSize"`
}

type SqliteCfg struct {
	DbPath   string `yaml:"DbPath"`
	DbName   string `yaml:"DbName"`
	Password string `yaml:"Password"`
}

func ReadConfig(config interface{}, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(config)
	if err != nil {
		return err
	}
	return nil
}

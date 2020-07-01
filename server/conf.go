package server

import (
	"gopkg.in/yaml.v2"
	"os"
)

var cfg Config

type Config struct {
	Http HttpCfg `yaml:"http"`
}

type HttpCfg struct {
	ListenPort       int  `yaml:"listen_port"`
	ReadTimeout      uint `yaml:"read_timeout"`
	ReadIdleTimeout  uint `yaml:"read_idle_timeout"`
	WriteTimeout     uint `yaml:"write_timeout"`
	WriteIdleTimeout uint `yaml:"write_idle_timeout"`
	MaxHeaderSize    int  `yaml:"max_header_size"`
}

func ReadConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

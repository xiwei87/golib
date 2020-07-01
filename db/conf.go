package db

import (
	"os"

	"gopkg.in/yaml.v2"
)

var DbCfg DbConfig

type DbConfig struct {
	Sqlite SqliteCfg `yaml:"sqlite"`
	Mysql  MysqlCfg  `yaml:"mysql"`
}

type SqliteCfg struct {
	DbPath string `yaml:"db_path"`
	DbName string `yaml:"db_name"`
}

type MysqlCfg struct {
	Server            string `yaml:"http_server"`
	ConnectionNum     int    `yaml:"conn_num"`
	ConnectionIdleNum int    `yaml:"conn_idle_num"`
	UserName          string `yaml:"user_name"`
	PassWord          string `yaml:"password"`
	DbName            string `yaml:"db_name"`
}

func ReadConfig(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(DbCfg)
	if err != nil {
		return err
	}
	return nil
}

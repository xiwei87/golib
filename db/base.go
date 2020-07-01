package db

import (
	"errors"
	"xorm.io/xorm"
)

type DB_TYPE int

const (
	MYSQL  DB_TYPE = 1
	SQLITE         = 2
)

//初始化数据库
func InitDb(dbType DB_TYPE, confPath string) (*xorm.Engine, error) {
	var (
		err error
	)
	if "" == confPath {
		return nil, errors.New("数据库配置文件为空")
	}
	if err = ReadConfig(confPath); err != nil {
		return nil, err
	}
	if dbType == SQLITE {
		return initSqliteDb()
	}
	return nil, errors.New("不支持该数据库类型")
}

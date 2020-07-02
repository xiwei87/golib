package db

import (
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type SqliteConfig struct {
	DbPath string `yaml:"db_path"`
	DbName string `yaml:"db_name"`
}

//创建本地目录
func createDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(dirName, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

//初始化SQLITE
func InitSqliteDb(cfg *SqliteConfig) (*xorm.Engine, error) {
	var (
		err    error
		engine *xorm.Engine
		dbPath string
	)
	//处理数据库目录地址
	dbPath = strings.TrimSuffix(cfg.DbPath, "/")
	//生成数据库目录
	if err = createDir(dbPath); err != nil {
		return nil, err
	}
	//创建数据库
	dbPath = dbPath + "/" + cfg.DbName + ".db"
	if engine, err = xorm.NewEngine("sqlite3", dbPath); err != nil {
		return nil, err
	}
	//设置时区
	engine.DatabaseTZ, _ = time.LoadLocation("Asia/Shanghai")

	return engine, nil
}

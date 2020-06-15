package db

import (
	_ "github.com/mattn/go-sqlite3"
	"gitlab.66ifuel.com/golang-tools/golib/config"
	"strings"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitSqliteDb() (err error) {
	var dbPath string
	dbPath = config.Cfg.Sqlite.DbPath
	strings.TrimSuffix(dbPath, "/")
	dbPath = dbPath + "/" + config.Cfg.Sqlite.DbName + ".db"
	if engine, err = xorm.NewEngine("sqlite3", dbPath); err != nil {
		return err
	}

	return nil
}

func QuerySqliteEngine() *xorm.Engine {
	return engine
}

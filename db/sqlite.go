package db

import (
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.66ifuel.com/golang-tools/golib/config"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func dbDirCreate(dbDir string) error {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(dbDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitSqliteDb() (err error) {
	var dbPath string

	dbPath = config.Cfg.Sqlite.DbPath
	strings.TrimSuffix(dbPath, "/")
	/* create db path */
	if err = dbDirCreate(dbPath); err != nil {
		return err
	}
	dbPath = dbPath + "/" + config.Cfg.Sqlite.DbName + ".db"
	if engine, err = xorm.NewEngine("sqlite3", dbPath); err != nil {
		return err
	}
	return nil
}

func QuerySqliteEngine() *xorm.Engine {
	return engine
}

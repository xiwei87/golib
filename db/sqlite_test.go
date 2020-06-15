package db

import (
	"gitlab.66ifuel.com/golang-tools/golib/config"
	"testing"
	"time"
	"xorm.io/xorm"
)

type User struct {
	Id        int64
	Name      string    `xorm:"varchar(25) notnull unique 'usr_name' comment('姓名')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (u *User) TableName() string {
	return "t_user"
}

func TestCreateSqliteDb(t *testing.T) {
	config.Cfg.Sqlite.DbName = "test"
	config.Cfg.Sqlite.DbPath = "./"

	if err := InitSqliteDb(); err != nil {
		t.Error(err)
	}
	var engine *xorm.Engine
	engine = QuerySqliteEngine()
	_ = engine.Sync2(new(User))
}

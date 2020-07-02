package db

import (
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
	var (
		err    error
		engine *xorm.Engine
		cfg    SqliteConfig
	)
	cfg.DbPath = "/tmp/test/"
	cfg.DbName = "user"

	if engine, err = InitSqliteDb(&cfg); err != nil {
		t.Error(err)
	}
	if err = engine.Sync2(new(User)); err != nil {
		t.Error(err)
	}
}

func TestInsetUser(t *testing.T) {
	var (
		err    error
		engine *xorm.Engine
		cfg    SqliteConfig
	)
	cfg.DbPath = "/tmp/test/"
	cfg.DbName = "user"

	if engine, err = InitSqliteDb(&cfg); err != nil {
		t.Error(err)
	}
	if err = engine.Sync2(new(User)); err != nil {
		t.Error(err)
	}
	userInfo := &User{
		Name: "lance1111",
	}
	if _, err = engine.Insert(userInfo); err != nil {
		t.Error(err)
	}
}

package dbs

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-wyvern/linear-chan/config"
)

var Db *gorm.DB

func InitMysql() error {
	var err error
	Db, err = gorm.Open("mysql", config.AppConfig.Mysql.GetAddress())
	if err != nil {
		return err
	}
	err = Db.DB().Ping()
	if err != nil {
		return err
	}
	Db.LogMode(config.AppConfig.Mysql.Debug)
	Db.DB().SetMaxIdleConns(config.AppConfig.Mysql.MaxIdle)
	Db.DB().SetMaxOpenConns(config.AppConfig.Mysql.MaxOpen)
	return nil
}

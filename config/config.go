package config

import (
	"fmt"

	"github.com/go-ini/ini"
)

var AppConfig =new(DeployConfig)
var ProjectConfig *ini.File

type DeployConfig struct {
	Mysql MysqlConfig `ini:"mysql"`
}

type MysqlConfig struct {
	UserName   string `ini:"username"`
	Password   string `ini:"password"`
	Database   string `ini:"database"`
	Address    string  `ini:"address"`
	Parameters string  `ini:"parameters"`
	Debug      bool `ini:"debug"`
	MaxIdle    int `ini:"max_idle"`
	MaxOpen    int `ini:"max_open"`
}

func (c MysqlConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s@%s/%s?%s",
		c.UserName, c.Password, c.Address, c.Database, c.Parameters)
}

func InitConfig() error {
	cfg, err := ini.Load("/etc/linear-chan.ini")
	if err != nil {
		return err
	}
	ProjectConfig, err = ini.Load("/etc/linear-chan.d/projects.ini")
	if err != nil {
		return err
	}
	err= cfg.MapTo(AppConfig)
	if err != nil {
		return err
	}
	return nil
}

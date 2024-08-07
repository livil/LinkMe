package ioc

import (
	"fmt"
	"github.com/GoSimplicity/LinkMe/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	DSN string `yaml:"dsn"`
}

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	var c config
	if err := viper.UnmarshalKey("db", &c); err != nil {
		panic(fmt.Errorf("init failed：%v", err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	// 初始化表
	if err = dao.InitTables(db); err != nil {
		panic(err)
	}
	return db
}

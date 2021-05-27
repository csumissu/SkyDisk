package model

import (
	"fmt"
	"time"

	"github.com/csumissu/SkyDisk/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DatabaseCfg.User,
		conf.DatabaseCfg.Password,
		conf.DatabaseCfg.Host,
		conf.DatabaseCfg.Port,
		conf.DatabaseCfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Can not connect to the database! %v", err))
	}

	db.Debug()
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	DB = db

	migration()
}

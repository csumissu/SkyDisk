package infra

import (
	"fmt"
	"github.com/csumissu/SkyDisk/config"
	"github.com/csumissu/SkyDisk/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var GormDB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DatabaseCfg.User,
		config.DatabaseCfg.Password,
		config.DatabaseCfg.Host,
		config.DatabaseCfg.Port,
		config.DatabaseCfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		util.Logger.Panic("can not connect to the database! %v", err)
	}

	db.Debug()
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	GormDB = db
}

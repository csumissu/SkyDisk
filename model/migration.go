package model

import (
	"github.com/csumissu/SkyDisk/infrastructure"
	"github.com/csumissu/SkyDisk/util/logger"
	"gorm.io/gorm"
)

var db = infrastructure.GormDB

func init() {
	if err := db.AutoMigrate(&User{}); err != nil {
		logger.Fatal("could not migrate schema, %v", err)
	}

	addDefaultUser()
}

func addDefaultUser() {
	_, err := GetUserByID(1)

	if err == gorm.ErrRecordNotFound {
		defaultUser := new(User)
		defaultUser.Nickname = "Evan"
		defaultUser.Username = "admin"
		defaultUser.Status = Active
		defaultUser.Password = "123"

		if err := defaultUser.Create(); err != nil {
			logger.Fatal("fail to create default user, %v", err)
		}
	}
}

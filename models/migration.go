package models

import (
	"github.com/csumissu/SkyDisk/infra"
	"github.com/csumissu/SkyDisk/util"
	"gorm.io/gorm"
)

var db = infra.GormDB

func init() {
	if err := db.AutoMigrate(&User{}, &Object{}); err != nil {
		util.Logger.Panic("could not migrate schema, %v", err)
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
			util.Logger.Panic("fail to create default user, %v", err)
		}
	}
}

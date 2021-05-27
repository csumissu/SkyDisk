package model

import (
	"fmt"
	"gorm.io/gorm"
)

func migration() {
	if err := DB.AutoMigrate(&User{}); err != nil {
		panic(fmt.Sprintf("Could not migrate schema, %v", err))
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
			panic(fmt.Sprintf("Fail to create default user, %v", err))
		}
	}
}

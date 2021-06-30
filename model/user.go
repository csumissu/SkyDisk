package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Nickname string `gorm:"not null"`
	Status   int    `gorm:"not null"`
}

const (
	Active = iota
	Banned
)

func GetUserByUsername(username string) (User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	return user, result.Error
}

func GetUserByID(ID interface{}) (User, error) {
	var user User
	result := db.First(&user, ID)
	return user, result.Error
}

func GetActiveUserByID(ID interface{}) (User, error) {
	var user User
	result := db.Where("status = ?", Active).First(&user, ID)
	return user, result.Error
}

func (user *User) CheckPassword(password string) (bool, error) {
	return user.Password == password, nil
}

func (user *User) Create() error {
	return db.Create(user).Error
}

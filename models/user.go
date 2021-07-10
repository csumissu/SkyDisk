package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at;index:idx_only_one,unique"`
	Username  string         `gorm:"not null;index:idx_user_name;index:idx_only_one,unique"`
	Password  string         `gorm:"not null"`
	Nickname  string         `gorm:"not null"`
	Status    Status         `gorm:"not null;index:idx_status"`
}

type Status int

const (
	Active Status = iota
	Banned
)

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	result := db.Where("username = ?", username).First(user)
	return user, result.Error
}

func GetUserByID(ID interface{}) (*User, error) {
	user := &User{}
	result := db.First(user, ID)
	return user, result.Error
}

func GetActiveUserByID(ID interface{}) (*User, error) {
	user := &User{}
	result := db.Where("status = ?", Active).First(user, ID)
	return user, result.Error
}

func (user User) CheckPassword(password string) (bool, error) {
	return user.Password == password, nil
}

func (user *User) Create() error {
	return db.Create(user).Error
}

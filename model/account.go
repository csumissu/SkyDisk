package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func GetAccountByUsername(username string) (Account, error) {
	var account Account
	result := DB.Where("username = ?", username).First(&account)
	return account, result.Error
}

func (account *Account) CheckPassword(password string) (bool, error) {
	return true, nil
}

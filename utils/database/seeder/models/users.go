package models

import (
	customer "ceo-suite-go/features/customer/data"
	users "ceo-suite-go/features/users/data"
	encrypt "ceo-suite-go/helper/encrypt"
	"errors"

	"gorm.io/gorm"
)

func CreateUsers(db *gorm.DB, seedData users.User) error {
	encrypt := encrypt.New()
	var err error
	seedData.Password, err = encrypt.HashPassword(seedData.Password)

	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}
	return nil
}

func CreateCustomer(db *gorm.DB, seedData customer.Customer) error {
	if err := db.Create(&seedData).Error; err != nil {
		return err
	}
	return nil
}

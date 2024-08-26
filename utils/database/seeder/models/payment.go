package models

import (
	booking "ceo-suite-go/features/booking/data"

	"gorm.io/gorm"
)

func CreatePaymentList(db *gorm.DB, name string) error {

	seedData := booking.PaymentList{
		Name: name,
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}
	return nil
}

// func CreateCustomer(db *gorm.DB, seedData customers.Customer) error {
// 	if err := db.Create(&seedData).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

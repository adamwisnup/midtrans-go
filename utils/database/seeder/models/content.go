package models

import (
	"ceo-suite-go/features/banner"
	"ceo-suite-go/features/contactus"
	"ceo-suite-go/features/faq"
	"ceo-suite-go/features/information"

	loremipsum "gopkg.in/loremipsum.v1"

	"gorm.io/gorm"
)

func CreateFAQ(db *gorm.DB) error {

	seedData := faq.Faq{
		Question: loremipsum.New().Words(5) + " ?",
		Answer:   loremipsum.New().Words(7),
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}

	return nil
}

func CreateBanner(db *gorm.DB) error {

	seedData := banner.Banner{
		Title:       loremipsum.New().Words(5),
		URL:         "https://www.commercialcafe.com/blog/wp-content/uploads/sites/10/shutterstock_670450231.jpg",
		Description: loremipsum.New().Words(10),
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}

	return nil
}

func CreateInformation(db *gorm.DB) error {

	seedData := information.Information{
		Name:        loremipsum.New().Words(3),
		Description: loremipsum.New().Words(10),
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}

	return nil
}

func CreateContactUs(db *gorm.DB) error {

	seedData := contactus.ContactUs{
		Name:    loremipsum.New().Words(3),
		Email:   "example@email.com",
		Subject: loremipsum.New().Words(7),
		Message: loremipsum.New().Words(10),
	}

	if err := db.Create(&seedData).Error; err != nil {
		return err
	}

	return nil
}

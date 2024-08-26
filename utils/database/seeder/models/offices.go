package models

import (
	offices "ceo-suite-go/features/office/data"

	"gorm.io/gorm"
)

func CreateOffice(db *gorm.DB, newData offices.Office, newPrice offices.Price, details offices.OfficeDetails, newCatalogues []offices.OfficeCatalogue) error {
	// Create the office
	if err := db.Create(&newData).Error; err != nil {
		return err
	}

	// Create the price
	newPrice.OfficeID = newData.ID
	if err := db.Create(&newPrice).Error; err != nil {
		return err
	}

	details.OfficeID = newData.ID
	if err := db.Create(&details).Error; err != nil {
		return err
	}

	// Create the catalogues
	for i := range newCatalogues {
		newCatalogues[i].OfficeID = newData.ID
	}
	if err := db.Create(&newCatalogues).Error; err != nil {
		return err
	}

	return nil
}

func CreateOfficeCategory(db *gorm.DB, name string) error {
	if err := db.Create(&offices.OfficeCategory{
		CategoryName: name,
	}).Error; err != nil {
		return err
	}

	return nil
}

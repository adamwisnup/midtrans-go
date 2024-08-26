package models

import (
	regions "ceo-suite-go/features/regions/data"

	"gorm.io/gorm"
)

func CreateRegions(db *gorm.DB, name string) error {
	return db.Create(&regions.Region{Name: name}).Error
}

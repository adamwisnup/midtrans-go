package database

import (
	DataBanner "ceo-suite-go/features/banner/data"
	DataBooking "ceo-suite-go/features/booking/data"
	DataContactUs "ceo-suite-go/features/contactus/data"
	DataCustomer "ceo-suite-go/features/customer/data"
	DataFaq "ceo-suite-go/features/faq/data"
	DataInformation "ceo-suite-go/features/information/data"
	DataOffice "ceo-suite-go/features/office/data"
	DataPromo "ceo-suite-go/features/promo/data"
	DataRegions "ceo-suite-go/features/regions/data"
	DataUser "ceo-suite-go/features/users/data"
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {

	// USER DATA MANAGEMENT \\
	db.AutoMigrate(DataUser.User{})
	// db.AutoMigrate(DataUser.UserResetPass{})
	// db.AutoMigrate(DataRegions.Region{})
	fmt.Println("[MIGRATION] Success creating user and regions")

	// OFFICE DATA MANAGEMENT \\
	// db.AutoMigrate(DataOffice.Office{})
	// db.AutoMigrate(DataOffice.Catering{})
	// db.AutoMigrate(DataOffice.OfficeDetails{})
	// db.AutoMigrate(DataOffice.OfficeCatalogue{})
	// db.AutoMigrate(DataOffice.OfficeReview{})
	// db.AutoMigrate(DataOffice.Price{})
	// db.AutoMigrate(DataOffice.CateringPrice{})
	// db.AutoMigrate(DataOffice.OfficeCategory{})
	// fmt.Println("[MIGRATION] Success creating office master data")

	// BOOKING DATA MANAGEMENT \\
	// db.AutoMigrate(DataBooking.Booking{})
	// db.AutoMigrate(DataBooking.CateringDetails{})
	// db.AutoMigrate(DataBooking.Payment{})
	// db.AutoMigrate(DataBooking.PaymentList{})
	// fmt.Println("[MIGRATION] Success creating booking management table")

	// CONTENT MANAGEMENT \\
	// db.AutoMigrate(DataInformation.Information{})
	// db.AutoMigrate(DataInformation.InformationContact{})
	// db.AutoMigrate(DataContactUs.ContactUs{})
	// db.AutoMigrate(DataCustomer.Customer{})
	// db.AutoMigrate(DataFaq.Faq{})
	// db.AutoMigrate(DataBanner.Banner{})
	// fmt.Println("[MIGRATION] Success creating content master data")

	// PROMO MANAGEMENT \\
	// db.AutoMigrate(DataPromo.Discount{})
	// db.AutoMigrate(DataPromo.DiscountOffice{})
	// db.AutoMigrate(DataPromo.FlashSale{})
	// db.AutoMigrate(DataPromo.FlashSaleProduct{})
	// fmt.Println("[MIGRATION] Success creating promo master data")
}

func MigrateWithDrop(db *gorm.DB) {
	// JUST RUN 1 TIME
	db.Exec("DROP SEQUENCE IF EXISTS invoice_seq;")
	db.Exec("DROP SCHEMA public CASCADE;")
	db.Exec("CREATE SCHEMA public;")
	db.Exec("CREATE SEQUENCE invoice_seq START WITH 1 INCREMENT BY 1 CACHE 10;")

	fmt.Println("[MIGRATION] Success dropping and creating schema")

	db.Exec("CREATE TYPE roles AS ENUM ('CUSTOMER', 'SUPERADMIN');")
	db.Exec("CREATE TYPE status AS ENUM ('ACTIVE', 'NOTACTIVE');")
	db.Exec("CREATE TYPE promotion_type AS ENUM ('direct_discount', 'buy_x_get_x');")
	fmt.Println("[MIGRATION] Success creating enum types for roles and status")

	// USER DATA MANAGEMENT \\
	db.AutoMigrate(DataUser.User{})
	db.AutoMigrate(DataUser.UserResetPass{})
	db.AutoMigrate(DataRegions.Region{})
	fmt.Println("[MIGRATION] Success creating user and regions")

	// OFFICE DATA MANAGEMENT \\
	db.AutoMigrate(DataOffice.Office{})
	db.AutoMigrate(DataOffice.OfficeDetails{})
	db.AutoMigrate(DataOffice.OfficeCatalogue{})
	db.AutoMigrate(DataOffice.OfficeReview{})
	db.AutoMigrate(DataOffice.Price{})
	db.AutoMigrate(DataOffice.OfficeCategory{})
	fmt.Println("[MIGRATION] Success creating office master data")

	// BOOKING DATA MANAGEMENT \\
	db.AutoMigrate(DataBooking.Booking{})
	db.AutoMigrate(DataBooking.Payment{})
	db.AutoMigrate(DataBooking.PaymentList{})
	fmt.Println("[MIGRATION] Success creating booking management table")

	// CONTENT MANAGEMENT \\
	db.AutoMigrate(DataInformation.Information{})
	db.AutoMigrate(DataContactUs.ContactUs{})
	db.AutoMigrate(DataCustomer.Customer{})
	db.AutoMigrate(DataFaq.Faq{})
	db.AutoMigrate(DataBanner.Banner{})
	fmt.Println("[MIGRATION] Success creating content master data")

	// PROMO MANAGEMENT \\
	db.AutoMigrate(DataPromo.Discount{})
	db.AutoMigrate(DataPromo.DiscountOffice{})
	db.AutoMigrate(DataPromo.FlashSale{})
	db.AutoMigrate(DataPromo.FlashSaleProduct{})
	fmt.Println("[MIGRATION] Success creating promo master data")

}

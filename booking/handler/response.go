package handler

import (
	booking "ceo-suite-go/features/booking"
	order "ceo-suite-go/features/booking"
	"time"

	"github.com/google/uuid"
)

type InputResponse struct {
	BookingID uuid.UUID `json:"id"`
	// BookingInvoice   string                    `json:"booking_invoice" form:"booking_invoice"`
	// CustomerName     string                    `json:"customer_name" form:"customer_name"`
	// HourCount        uint                      `json:"hour_count" form:"hour_count"`
	// DaysCount        uint                      `json:"days_count" form:"days_count"`
	// HalfDayCount     uint                      `json:"half_day_count"`
	// WeekCount        uint                      `json:"week_count" form:"week_count"`
	// MonthCount       uint                      `json:"month_count" form:"month_count"`
	// PaymentType      string                    `json:"payment_type" form:"payment_type"`
	// BookingPrice     uint                      `json:"booking_price"`
	// Discount         uint                      `json:"discount" form:"discount"`
	// OvertimeCount    uint                      `json:"overtime_hour_count"`
	// OvertimePrice    uint                      `json:"overtime_price"`
	// Catering         []booking.CateringDetails `json:"catering_details"`
	// CateringPrice    uint                      `json:"catering_total_price"`
	// Deposit          uint                      `json:"deposit"`
	// TotalPrice       uint                      `json:"total_price" form:"total_price"`
	// PPN              uint                      `json:"ppn"`
	// FinalPrice       uint                      `json:"final_price"`
	// BookingStatus    uint                      `json:"booking_status" form:"booking_status"`
	// BookingStartDate time.Time                 `json:"booking_start_date" form:"booking_start_date"`
	// BookingEndDate   time.Time                 `json:"booking_end_date" form:"booking_end_date"`
	// UsesDescription  string                    `json:"uses_description" form:"uses_description"`
	// Notes            string                    `json:"notes"`
}

type UpdateResponse struct {
	BookingID        uuid.UUID                 `json:"id"`
	BookingInvoice   string                    `json:"booking_invoice" form:"booking_invoice"`
	CustomerName     string                    `json:"customer_name" form:"customer_name"`
	HourCount        uint                      `json:"hour_count" form:"hour_count"`
	DaysCount        uint                      `json:"days_count" form:"days_count"`
	HalfDayCount     uint                      `json:"half_day_count"`
	WeekCount        uint                      `json:"week_count" form:"week_count"`
	MonthCount       uint                      `json:"month_count" form:"month_count"`
	PaymentType      string                    `json:"payment_type" form:"payment_type"`
	BookingPrice     uint                      `json:"booking_price"`
	Discount         uint                      `json:"discount" form:"discount"`
	OvertimeCount    uint                      `json:"overtime_hour_count"`
	OvertimePrice    uint                      `json:"overtime_price"`
	Catering         []booking.CateringDetails `json:"catering_details"`
	CateringPrice    uint                      `json:"catering_total_price"`
	Deposit          uint                      `json:"deposit"`
	TotalPrice       uint                      `json:"total_price" form:"total_price"`
	PPN              uint                      `json:"ppn"`
	FinalPrice       uint                      `json:"final_price"`
	BookingStatus    uint                      `json:"booking_status" form:"booking_status"`
	BookingStartDate time.Time                 `json:"booking_start_date" form:"booking_start_date"`
	BookingEndDate   time.Time                 `json:"booking_end_date" form:"booking_end_date"`
	UsesDescription  string                    `json:"uses_description" form:"uses_description"`
	Notes            string                    `json:"notes"`
}

type BookingGuestResponse struct {
	IsGuestMode bool      `json:"is_guest_mode"`
	Customer    any       `json:"customer"`
	BookingID   uuid.UUID `json:"booking_id"`
}

type GetResponse struct {
	MainData   any   `json:"data"`
	Page       *uint `json:"page"`
	PageSize   *uint `json:"page_size"`
	TotalPage  uint  `json:"total_page"`
	TotalItems uint  `json:"total_items"`
}

type BookingReport struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id"`
	OfficeName       string         `json:"office_name"`
	CustomerName     string         `json:"customer_name"`
	OfficeID         uint           `json:"office_id" gorm:"column:office_id"`
	UserID           uint           `json:"user_id,omitempty" gorm:"column:user_id"`
	CustomerID       uint           `json:"customer_id" gorm:"column:customer_id"`
	TotalPrice       uint           `json:"total_price" gorm:"column:total_price"`
	HourCount        uint           `json:"hour_count" gorm:"column:hour_count"`
	DaysCount        uint           `json:"days_count" gorm:"column:days_count"`
	HalfDayCount     uint           `json:"half_day_count" gorm:"column:half_day_count"`
	WeekCount        uint           `json:"week_count" gorm:"column:week_count"`
	MonthCount       uint           `json:"month_count" gorm:"column:month_count"`
	PaymentType      string         `json:"payment_type" gorm:"column:payment_type"`
	Discount         uint           `json:"discount" gorm:"column:discount"`
	BookingStatus    uint           `json:"booking_status" gorm:"column:booking_status"`
	BookingStartDate time.Time      `json:"booking_start_date" gorm:"column:booking_start_date"`
	BookingEndDate   time.Time      `json:"booking_end_date" gorm:"column:booking_end_date"`
	UsesDescription  string         `json:"uses_description" gorm:"column:uses_description"`
	VAAccount        string         `json:"va_account" gorm:"column:va_account"`
	Notes            string         `json:"notes" gorm:"column:notes"`
	Payment          *order.Payment `json:"payment,omitempty" gorm:"foreignKey:BookingID"`
}

package data

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BOOKING STATUS
// Please check again.
const (
	BookingCreated           = 1 /* Status when booking is created */
	BookingWaitingForPayment = 2 /* Status when Waiting For Payment */
	BookingSuccess           = 3 /* Status when Success */
	BookingFailed            = 4 /* Status when Failed */

	// BookingExpired            = 9  /* Status when Expired */
	BookingCanceledByCustomer = 10 /* Status when Canceled By Customer */
	BookingCanceledByMidtrans = 11 /* Status when Canceled By Midtrans */
	BookingExpired            = 12 /* Status when Expired */
	BookingCanceledByAdmin    = 13 /* Status when Canceled By Admin */

)

// MIDTRANS STATUS
const (
	PaymentChallenge = 1 /* Status when Payment is Challenged */
	PaymentSuccess   = 2 /* Status when Payment is Success */
	PaymentDenied    = 3 /* Status when Payment is Denied */
	PaymentCanceled  = 4 /* Status when Payment is Canceled */
	PaymentPending   = 5 /* Status when Payment is Pending */
)

type Booking struct {
	*gorm.Model
	ID                    uuid.UUID         `gorm:"column:id;type:uuid"`
	BookingInvoice        string            `gorm:"column:booking_invoice;type:varchar(255)"`
	OfficeID              uint              `gorm:"column:office_id"`
	CustomerID            uint              `gorm:"column:customer_id"`
	UserID                uint              `gorm:"column:user_id;default:0"`
	HourCount             uint              `gorm:"column:hour_count"`
	DaysCount             uint              `gorm:"column:days_count"`
	WeekCount             uint              `gorm:"column:week_count"`
	MonthCount            uint              `gorm:"column:month_count"`
	PaymentType           string            `gorm:"column:payment_type;type:varchar(255)"`
	BookingPrice          uint              `gorm:"column:booking_price"`
	OvertimeCount         uint              `gorm:"column:overtime_count"`
	OvertimePrice         uint              `gorm:"column:overtime_price"`
	CateringDetails       []CateringDetails `gorm:"foreignKey:BookingID"`
	CateringPrice         uint              `gorm:"column:catering_price"`
	TotalPrice            uint              `gorm:"column:total_price"`
	PPN                   uint              `gorm:"column:ppn"`
	Deposit               uint              `gorm:"column:deposit;default:0"`
	FinalPrice            uint              `gorm:"column:final_price"`
	Discount              uint              `gorm:"column:discount"`
	BookingStatus         uint              `gorm:"column:booking_status"`
	BookingStartDate      time.Time         `gorm:"column:booking_start_date"`
	BookingEndDate        time.Time         `gorm:"column:booking_end_date"`
	BookingExpirationTime time.Time         `gorm:"column:booking_expiration_time"`
	VAAccount             string            `gorm:"column:va_account;type:varchar(255)"`
	UsesDescription       string            `gorm:"column:uses_description"`
	Notes                 string            `gorm:"column:notes;default:null"`
	ReasonCancel          string            `gorm:"column:reason_cancel;type:text;default:null"`
	Payment               Payment           `gorm:"foreignKey:BookingID"`
}

type CateringDetails struct {
	*gorm.Model
	BookingID        uuid.UUID `gorm:"column:booking_id"`
	CateringID       uint      `gorm:"column:catering_id;default:0"`
	CateringCountPax uint      `gorm:"column:catering_count_pax;default:0"`
	CateringPrice    uint      `gorm:"column:catering_price;default:0"`
}

type Payment struct {
	*gorm.Model
	ID             uuid.UUID `gorm:"column:id;type:uuid"`
	BookingID      uuid.UUID `gorm:"column:booking_id"`
	PaymentStatus  uint      `gorm:"column:payment_status"`
	PaymentType    string    `gorm:"column:payment_type;type:varchar(255)"`
	TotalPrice     uint      `gorm:"column:total_price"`
	PaymentInvoice string    `gorm:"column:payment_invoice;type:varchar(255)"`
}

type PaymentList struct {
	*gorm.Model
	Name string `gorm:"column:name;type:varchar(255)"`
}

func (Booking) TableName() string {
	return "bookings"
}

func (Payment) TableName() string {
	return "payments"
}

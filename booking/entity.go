package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookingDemo struct {
	ID               uuid.UUID `json:"id" gorm:"column:id"`
	BookingInvoice   string    `json:"booking_invoice" gorm:"column:booking_invoice"`
	OfficeName       string    `json:"office_name"`
	CustomerName     string    `json:"customer_name"`
	OvertimePrice    uint      `json:"overtime_price" gorm:"column:overtime_price"`
	OvertimeCount    uint      `json:"overtime_count" gorm:"column:overtime_count"`
	TotalPrice       uint      `json:"total_price" gorm:"column:total_price"`
	PaymentType      string    `json:"payment_type" gorm:"column:payment_type"`
	Discount         uint      `json:"discount" gorm:"column:discount"`
	BookingStatus    uint      `json:"booking_status" gorm:"column:booking_status"`
	BookingStartDate time.Time `json:"booking_start_date" gorm:"column:booking_start_date"`
	BookingEndDate   time.Time `json:"booking_end_date" gorm:"column:booking_end_date"`
}

type Booking struct {
	ID                    uuid.UUID         `json:"id" gorm:"column:id"`
	BookingInvoice        string            `json:"booking_invoice" gorm:"column:booking_invoice"`
	OfficeName            string            `json:"office_name"`
	CustomerName          string            `json:"customer_name" gorm:"->"`
	OfficeID              uint              `json:"office_id" gorm:"column:office_id"`
	UserID                uint              `json:"user_id,omitempty" gorm:"column:user_id"`
	CustomerID            uint              `json:"customer_id" gorm:"column:customer_id"`
	HourCount             uint              `json:"hour_count" gorm:"column:hour_count"`
	DaysCount             uint              `json:"days_count" gorm:"column:days_count"`
	WeekCount             uint              `json:"week_count" gorm:"column:week_count"`
	MonthCount            uint              `json:"month_count" gorm:"column:month_count"`
	PaymentType           string            `json:"payment_type" gorm:"column:payment_type"`
	BookingPrice          uint              `json:"booking_price" gorm:"column:booking_price"`
	OvertimeCount         uint              `json:"overtime_count" gorm:"column:overtime_count"`
	OvertimePrice         uint              `json:"overtime_price" gorm:"column:overtime_price"`
	CateringDetails       []CateringDetails `json:"catering_details" gorm:"foreignKey:BookingID"`
	CateringPrice         uint              `json:"catering_price" gorm:"column:catering_price"`
	TotalPrice            uint              `json:"total_price" gorm:"column:total_price"`
	PPN                   uint              `json:"ppn" gorm:"column:ppn"`
	Deposit               uint              `json:"deposit" gorm:"column:deposit"`
	FinalPrice            uint              `json:"final_price" gorm:"final_price"`
	Discount              uint              `json:"discount" gorm:"column:discount"`
	BookingStatus         uint              `json:"booking_status" gorm:"column:booking_status"`
	BookingStartDate      time.Time         `json:"booking_start_date" gorm:"column:booking_start_date"`
	BookingEndDate        time.Time         `json:"booking_end_date" gorm:"column:booking_end_date"`
	BookingExpirationTime time.Time         `json:"booking_expiration_time" gorm:"column:booking_expiration_time"`
	UsesDescription       string            `json:"uses_description" gorm:"column:uses_description"`
	VAAccount             string            `json:"va_account" gorm:"column:va_account"`
	Notes                 string            `json:"notes" gorm:"column:notes"`
	ReasonCancel          string            `json:"reason_cancel" gorm:"column:reason_cancel"`
	Payment               *Payment          `json:"payment,omitempty" gorm:"foreignKey:BookingID"`
}

type CateringDetails struct {
	BookingID         uuid.UUID `json:"booking_id" gorm:"column:booking_id"`
	CateringID        uint      `json:"catering_id" gorm:"column:catering_id"`
	CateringName      string    `json:"catering_name" gorm:"-"`
	CateringUnitPrice *uint     `json:"catering_unit_price" gorm:"-"`
	CateringCountPax  uint      `json:"catering_count_pax" gorm:"column:catering_count_pax"`
	CateringPrice     uint      `json:"catering_price" gorm:"column:catering_price"`
}

type BookingCustomer struct {
	ID                    uuid.UUID `json:"id" gorm:"column:id"`
	BookingInvoice        string    `json:"booking_invoice" gorm:"column:booking_invoice"`
	UserID                uint      `json:"user_id,omitempty" gorm:"column:user_id"`
	CustomerID            uint      `json:"customer_id" gorm:"column:customer_id"`
	CustomerName          string    `json:"customer_name" gorm:"->"`
	OfficeID              uint      `json:"office_id" gorm:"column:office_id"`
	OfficeName            string    `json:"office_name"`
	OvertimePrice         uint      `json:"overtime_price" gorm:"column:overtime_price"`
	OvertimeCount         uint      `json:"overtime_count" gorm:"column:overtime_count"`
	CateringPrice         uint      `json:"catering_price" gorm:"column:catering_price"`
	FinalPrice            uint      `json:"final_price" gorm:"column:final_price"`
	PaymentType           string    `json:"payment_type" gorm:"column:payment_type"`
	Discount              uint      `json:"discount" gorm:"column:discount"`
	BookingStatus         uint      `json:"booking_status" gorm:"column:booking_status"`
	BookingStartDate      time.Time `json:"booking_start_date" gorm:"column:booking_start_date"`
	BookingEndDate        time.Time `json:"booking_end_date" gorm:"column:booking_end_date"`
	BookingExpirationTime time.Time `json:"booking_expiration_time" gorm:"column:booking_expiration_time"`
	VAAccount             string    `json:"va_account" gorm:"column:va_account"`
	UsesDescription       string    `json:"uses_description" gorm:"column:uses_description"`
	ReasonCancel          string    `json:"reason_cancel" gorm:"column:reason_cancel"`
	CreatedAt             time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type BookingCustomerDetails struct {
	ID                    uuid.UUID         `json:"id" gorm:"column:id"`
	BookingInvoice        string            `json:"booking_invoice" gorm:"column:booking_invoice"`
	OfficeID              uint              `json:"office_id" gorm:"column:office_id"`
	UserID                uint              `json:"user_id,omitempty" gorm:"column:user_id"`
	CustomerID            uint              `json:"customer_id" gorm:"column:customer_id"`
	CustomerName          string            `json:"customer_name" gorm:"->"`
	OfficeName            string            `json:"office_name"`
	HourCount             uint              `json:"hour_count" gorm:"column:hour_count"`
	DaysCount             uint              `json:"days_count" gorm:"column:days_count"`
	WeekCount             uint              `json:"week_count" gorm:"column:week_count"`
	MonthCount            uint              `json:"month_count" gorm:"column:month_count"`
	PaymentType           string            `json:"payment_type" gorm:"column:payment_type"`
	OvertimeCount         uint              `json:"overtime_count" gorm:"column:overtime_count"`
	BookingPrice          uint              `json:"booking_price" gorm:"column:booking_price"`
	OvertimePrice         uint              `json:"overtime_price" gorm:"column:overtime_price"`
	CateringDetails       []CateringDetails `json:"catering_details" gorm:"foreignKey:BookingID"`
	CateringPrice         uint              `json:"catering_price" gorm:"column:catering_price"`
	TotalPrice            uint              `json:"total_price" gorm:"column:total_price"`
	PPN                   uint              `json:"ppn" gorm:"column:ppn"`
	Deposit               uint              `json:"deposit" gorm:"column:deposit"`
	FinalPrice            uint              `json:"final_price" gorm:"column:final_price"`
	Discount              uint              `json:"discount" gorm:"column:discount"`
	BookingStatus         uint              `json:"booking_status" gorm:"column:booking_status"`
	BookingStartDate      time.Time         `json:"booking_start_date" gorm:"column:booking_start_date"`
	BookingEndDate        time.Time         `json:"booking_end_date" gorm:"column:booking_end_date"`
	BookingExpirationTime time.Time         `json:"booking_expiration_time" gorm:"column:booking_expiration_time"`
	UsesDescription       string            `json:"uses_description" gorm:"column:uses_description"`
	VAAccount             string            `json:"va_account" gorm:"column:va_account"`
	Notes                 string            `json:"notes" gorm:"column:notes"`
	ReasonCancel          string            `json:"reason_cancel" gorm:"column:reason_cancel"`
	CreatedAt             time.Time         `json:"created_at" gorm:"column:created_at"`
	UpdatedAt             time.Time         `json:"updated_at" gorm:"column:updated_at"`
	Payment               *Payment          `json:"payment" gorm:"foreignKey:BookingID"`
}

type Payment struct {
	BookingID      uuid.UUID `json:"booking_id" gorm:"column:booking_id"`
	PaymentStatus  uint      `json:"payment_status" gorm:"column:payment_status"`
	PaymentType    string    `json:"payment_type" gorm:"column:payment_type"`
	TotalPrice     uint      `json:"total_price" gorm:"column:total_price"`
	PaymentInvoice string    `json:"payment_invoice" gorm:"column:payment_invoice"`
	ExpiredAt      time.Time `json:"expired_at" gorm:"-"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type BookingGuestRequest struct {
	FullName    string `json:"full_name" form:"full_name"`
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type PaymentList struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"column:name"`
}

type BookingHandlerInterface interface {
	Checkout() echo.HandlerFunc
	CheckoutWithSnap() echo.HandlerFunc
	CreateBooking() echo.HandlerFunc
	GetAllBooking() echo.HandlerFunc
	GetAllBookingDemo() echo.HandlerFunc
	GetBookingByID() echo.HandlerFunc
	GetBookingByBookingInvoice() echo.HandlerFunc
	UpdateBooking() echo.HandlerFunc
	DeleteBooking() echo.HandlerFunc
	NotifBooking() echo.HandlerFunc
	GetAllPaymentList() echo.HandlerFunc
	GetPaymentByID() echo.HandlerFunc
	GetCustomersBooking() echo.HandlerFunc
	GetCustomersBookingByID() echo.HandlerFunc
	CreateBookingGuest() echo.HandlerFunc
	ReportBooking() echo.HandlerFunc
	CancelBooking() echo.HandlerFunc
	CancelBookingByAdmin() echo.HandlerFunc
	// CheckBookingAvailability() echo.HandlerFunc
}

type BookingServiceInterface interface {
	GetAllPaymentList() ([]PaymentList, error)
	GetPaymentByID(id uint) (*PaymentList, error)
	Checkout(bookingID uuid.UUID) (*Payment, interface{}, error)
	CreateBooking(newData Booking) (*Booking, interface{}, *string, error)
	CreateCateringDetails(newData CateringDetails) (*CateringDetails, error)
	GetAllBooking(search string, page uint, pageSize uint, status uint) ([]Booking, uint, uint, error)
	GetAllBookingDemo(search string, page uint, pageSize uint, status uint) ([]BookingDemo, uint, uint, error)
	GetBookingByID(id uuid.UUID) (*Booking, error)
	GetBookingByBookingInvoice(bookingInvoice string) (*Booking, error)
	UpdateBooking(id uuid.UUID, newData Booking) (*Booking, *string, error)
	NotifBooking(notificationPayload map[string]interface{}, newData Payment) (bool, error)
	DeleteBooking(id uuid.UUID) (bool, error)
	GetCustomersBooking(userID uint, status uint) ([]BookingCustomer, error)
	GetCustomersBookingByID(id uuid.UUID) (*BookingCustomerDetails, error)
	CreateBookingGuest(newData Booking, customer BookingGuestRequest) (*Booking, *BookingGuestRequest, error)
	ReportBooking() ([]Booking, error)
	CheckBookingAvailability(officeID uint, startTime time.Time, endTime time.Time) (bool, error)
	CheckoutWithSnap(bookingID uuid.UUID) (*Payment, interface{}, error)
	CancelBooking(bookingID uuid.UUID, explanation string) (bool, error)
	CalculateDuration(start, end time.Time) (normalDuration, overtimeDuration *time.Duration, err error)
}

type BookingDataInterface interface {
	GetAllPaymentList() ([]PaymentList, error)
	GetPaymentByID(id uint) (*PaymentList, error)
	Checkout(bookingID uuid.UUID) (*Payment, error)
	CreateBooking(newData Booking) (*Booking, *Payment, *string, error)
	CreateCateringDetails(newData CateringDetails) (*CateringDetails, error)
	GetAllBooking(search string, page uint, pageSize uint, status uint) ([]Booking, uint, uint, error)
	GetAllBookingDemo(search string, page uint, pageSize uint, status uint) ([]BookingDemo, uint, uint, error)
	GetBookingByID(id uuid.UUID) (*Booking, error)
	GetBookingByBookingInvoice(bookingInvoice string) (*Booking, error)
	GetAndUpdatePayment(newData Payment, id string) (bool, error)
	GetPaymentDataByID(invoice string) (*Payment, error)
	UpdateBooking(id uuid.UUID, newData Booking) (*Booking, *string, error)
	DeleteBooking(id uuid.UUID) (bool, error)
	UpdateBookingStatus(id uuid.UUID, newData Booking) (bool, error)
	UpdatePaymentByBookingID(bookingID uuid.UUID, newData Payment) (bool, error)
	GetCustomersBooking(userID uint, status uint) ([]BookingCustomer, error)
	GetCustomersBookingByID(id uuid.UUID) (*BookingCustomerDetails, error)
	CreateBookingGuest(newData Booking, custID uint) (*Booking, error)
	ReportBooking() ([]Booking, error)
	CheckBookingAvailability(officeID uint, startTime time.Time, endTime time.Time) (bool, error)
	CancelBooking(bookingID uuid.UUID, explanation string) (bool, error)
	CalculateDuration(start, end time.Time) (normalDuration, overtimeDuration time.Duration, err error)
	UpdateExpiredBookings() (int64, error)
}

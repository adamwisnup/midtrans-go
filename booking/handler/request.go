package handler

import (
	"time"

	"github.com/google/uuid"
)

type BookingGuestRequest struct {
	FullName    string `json:"full_name" form:"full_name"`
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`

	OfficeID              uint              `json:"office_id" form:"office_id"`
	CateringRequest       []CateringRequest `json:"catering"`
	UserID                uint              `json:"user_id" form:"user_id"`
	BookingStartDate      time.Time         `json:"booking_start_date" form:"booking_start_date"`
	BookingEndDate        time.Time         `json:"booking_end_date" form:"booking_end_date"`
	UseOvertimeAC         bool              `json:"use_overtime_ac" form:"use_overtime_ac"`
	AcceptTermsConditions bool              `json:"accept_terms_conditions" form:"accept_terms_conditions"`
	UsesDescription       string            `json:"uses_description" form:"uses_description"`
}
type InputRequest struct {
	OfficeID         uint      `json:"office_id" form:"office_id"`
	CateringID       []uint    `json:"catering_id" form:"catering_id"`
	CateringPaxCount []uint    `json:"catering_pax_count" form:"catering_pax_count"`
	UserID           uint      `json:"user_id" form:"user_id"`
	BookingStartDate time.Time `json:"booking_start_date" form:"booking_start_date"`
	BookingEndDate   time.Time `json:"booking_end_date" form:"booking_end_date"`
	UsesDescription  string    `json:"uses_description" form:"uses_description"`
}

type InputRequestTest struct {
	OfficeID              uint              `json:"office_id" form:"office_id"`
	CateringRequest       []CateringRequest `json:"catering"`
	UserID                uint              `json:"user_id" form:"user_id"`
	BookingStartDate      time.Time         `json:"booking_start_date" form:"booking_start_date"`
	BookingEndDate        time.Time         `json:"booking_end_date" form:"booking_end_date"`
	UseOvertimeAC         bool              `json:"use_overtime_ac" form:"use_overtime_ac"`
	AcceptTermsConditions bool              `json:"accept_terms_conditions" form:"accept_terms_conditions"`
	UsesDescription       string            `json:"uses_description" form:"uses_description"`
}

type CancelRequest struct {
	Explanation string `json:"explanation" form:"explanation"`
}

type CateringRequest struct {
	CateringID       uint `json:"catering_id" form:"catering_id"`
	CateringPaxCount uint `json:"catering_pax_count" form:"catering_pax_count"`
}

type CheckoutRequest struct {
	BookingID uuid.UUID `json:"booking_id" form:"booking_id"`
}

type UpdateRequest struct {
	OfficeID         uint      `json:"office_id" form:"office_id"`
	UserID           uint      `json:"user_id" form:"user_id"`
	PaymentType      string    `json:"payment_type" form:"payment_type"`
	BookingStartDate time.Time `json:"booking_start_date" form:"booking_start_date"`
	BookingEndDate   time.Time `json:"booking_end_date" form:"booking_end_date"`
	UsesDescription  string    `json:"uses_description" form:"uses_description"`
}

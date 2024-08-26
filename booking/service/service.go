package service

import (
	booking "ceo-suite-go/features/booking"
	"ceo-suite-go/features/customer"
	"ceo-suite-go/features/users"
	"ceo-suite-go/helper"
	"ceo-suite-go/helper/email"
	mt "ceo-suite-go/utils/midtrans"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
)

type BookingService struct {
	cs customer.CustomerDataInterface
	us users.UserServiceInterface
	d  booking.BookingDataInterface
	j  helper.JWTInterface
	em email.EmailInterface
	mt mt.MidtransService
}

func New(data booking.BookingDataInterface, us users.UserServiceInterface, jwt helper.JWTInterface, em email.EmailInterface, mt mt.MidtransService, cs customer.CustomerDataInterface) booking.BookingServiceInterface {
	return &BookingService{
		d:  data,
		us: us,
		j:  jwt,
		em: em,
		mt: mt,
		cs: cs,
	}
}

func (bs *BookingService) GetAllPaymentList() ([]booking.PaymentList, error) {
	res, err := bs.d.GetAllPaymentList()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bs *BookingService) GetPaymentByID(id uint) (*booking.PaymentList, error) {
	res, err := bs.d.GetPaymentByID(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bs *BookingService) GetCustomersBooking(userID uint, status uint) ([]booking.BookingCustomer, error) {
	res, err := bs.d.GetCustomersBooking(userID, status)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bs *BookingService) GetCustomersBookingByID(id uuid.UUID) (*booking.BookingCustomerDetails, error) {
	res, err := bs.d.GetCustomersBookingByID(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (bs *BookingService) NotifBooking(notificationPayload map[string]interface{}, newData booking.Payment) (bool, error) {
	paymentStatus, invoice, err := bs.mt.TransactionStatus(notificationPayload)
	if err != nil {
		return false, err
	}

	newData.PaymentStatus = uint(paymentStatus)
	result, err := bs.d.GetAndUpdatePayment(newData, invoice)

	if err != nil {
		return false, err
	}

	fmt.Println("Payment : ", newData)

	if newData.PaymentStatus == 2 {

		payment, err := bs.d.GetPaymentDataByID(invoice)

		if err != nil {
			return false, err
		}

		resBook, err := bs.d.GetBookingByID(payment.BookingID)
		if err != nil {
			return false, err
		}

		user, err := bs.us.GetProfile(int(resBook.UserID))
		if err != nil {
			return false, err
		}

		header, body := bs.em.HTMLBodyBookingSuccess(*resBook)

		err = bs.em.SendEmail(user.Email, header, body)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}

func (bs *BookingService) CreateBookingGuest(newData booking.Booking, custData booking.BookingGuestRequest) (*booking.Booking, *booking.BookingGuestRequest, error) {
	_, err := bs.CheckBookingAvailability(newData.OfficeID, newData.BookingStartDate, newData.BookingEndDate)

	if err != nil {
		return nil, nil, err
	}

	password := helper.RandomString(7)

	user, err := bs.us.Register(users.User{
		Name:            custData.FullName,
		Email:           custData.Email,
		Password:        password,
		PhoneNumber:     custData.PhoneNumber,
		Role:            "CUSTOMER",
		IsEmailVerified: false,
	})

	if err != nil {
		return nil, nil, err
	}

	err = bs.us.ResendCodeVerifyEmailandGetPassword(custData.Email, password)

	if err != nil {
		return nil, nil, errors.New("send email error")
	}

	cust, err := bs.cs.CreateCustomer(customer.Customer{
		UserID:         user.ID,
		FullName:       custData.FullName,
		Position:       fmt.Sprintf("%v Position", custData.FullName),
		CompanyName:    fmt.Sprintf("%v Company Name", custData.FullName),
		CompanyEmail:   custData.Email,
		CompanyAddress: fmt.Sprintf("%v Company Address", custData.FullName),
	})

	if err != nil {
		return nil, nil, err
	}

	newData.PaymentType = "payment_gateway"
	newData.UserID = user.ID

	result, err := bs.d.CreateBookingGuest(newData, cust.ID)
	if err != nil {
		return nil, nil, err
	}

	return result, &custData, nil
}

func (bs *BookingService) CreateBooking(newData booking.Booking) (*booking.Booking, interface{}, *string, error) {
	newData.PaymentType = "payment_gateway"

	result, _, name, err := bs.d.CreateBooking(newData)
	if err != nil {
		return nil, nil, nil, err
	}

	return result, nil, name, nil
}

func (bs *BookingService) CreateCateringDetails(newData booking.CateringDetails) (*booking.CateringDetails, error) {
	result, err := bs.d.CreateCateringDetails(newData)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (bs *BookingService) Checkout(bookingID uuid.UUID) (*booking.Payment, interface{}, error) {
	res, err := bs.d.Checkout(bookingID)

	if err != nil {
		return nil, nil, err
	}

	// fmt.Println(res)

	bookDetail, _ := bs.d.GetBookingByID(bookingID)
	cust, _ := bs.cs.GetCustomerByID(bookDetail.CustomerID)

	custDetails := &midtrans.CustomerDetails{
		FName: cust.FullName,
		Email: cust.CompanyEmail,
	}

	itemDetails := []midtrans.ItemDetails{
		{
			ID:    fmt.Sprint(bookDetail.OfficeID),
			Name:  bookDetail.OfficeName,
			Price: int64(res.TotalPrice),
			Qty:   1,
		},
	}

	_, resultMidtrans, err := bs.mt.GenerateTransaction(int(res.TotalPrice), strings.ToLower(res.PaymentType), res.PaymentInvoice, custDetails, itemDetails)
	// fmt.Println("Resultmidtrans: ", resultMidtrans)
	if err != nil {
		return nil, nil, err
	}

	if resultMidtrans == nil {
		fmt.Println("Result midtrans : ", resultMidtrans)
		return nil, nil, errors.New("midtrans error, no response payload")
	}

	var vaAccount, callbackURL, valueToUpdate string

	vaAccountValue, ok := resultMidtrans["va_account"].(string)
	if ok {
		vaAccount = vaAccountValue
	}

	callbackURLValue, ok := resultMidtrans["callback_url"].(string)
	if ok {
		callbackURL = callbackURLValue
	}

	if vaAccount != "" {
		valueToUpdate = vaAccount
	} else {
		valueToUpdate = callbackURL
	}

	_, err = bs.d.UpdateBookingStatus(bookingID, booking.Booking{
		BookingStatus: 2,
		VAAccount:     valueToUpdate,
	})

	if err != nil {
		return nil, nil, err
	}

	return res, resultMidtrans, nil
}
func (bs *BookingService) CalculateDuration(start, end time.Time) (normalDuration, overtimeDuration *time.Duration, err error) {

	nd, od, err := bs.d.CalculateDuration(start, end)

	if err != nil {
		return nil, nil, err
	}

	return &nd, &od, nil
}

func (bs *BookingService) CheckoutWithSnap(bookingID uuid.UUID) (*booking.Payment, interface{}, error) {
	bookDetail, _ := bs.d.GetBookingByID(bookingID)

	if bookDetail.BookingStatus >= 2 {
		mt := make(map[string]string)
		mt["redirect_url"] = bookDetail.VAAccount
		return nil, mt, nil
	}

	if time.Now().Local().After(bookDetail.BookingExpirationTime) {
		return nil, nil, errors.New("this booking is already expired")
	}

	res, err := bs.d.Checkout(bookingID)

	if err != nil {
		return nil, nil, err
	}

	// fmt.Println(res)

	cust, _ := bs.cs.GetCustomerByID(bookDetail.CustomerID)

	custDetails := &midtrans.CustomerDetails{
		FName: cust.FullName,
		Email: cust.User.Email,
	}

	var bookOfficeName string = ""
	bookOfficeName = bookDetail.OfficeName
	bookSplit := strings.Split(bookDetail.OfficeName, " ")
	if len(bookSplit) > 4 {
		bookOfficeName = fmt.Sprintf("%v...", bookOfficeName[:16])
	}

	itemDetails := []midtrans.ItemDetails{
		{
			ID:    fmt.Sprint(bookDetail.OfficeID),
			Name:  fmt.Sprintf("%v incl. PPN 11%%", bookOfficeName),
			Price: int64(res.TotalPrice),
			Qty:   1,
		},
	}

	resultMidtrans, err := bs.mt.GenerateTransactionSnap(int(res.TotalPrice), strings.ToLower(res.PaymentType), res.PaymentInvoice, custDetails, itemDetails)
	// fmt.Println("Resultmidtrans: ", resultMidtrans)
	if err != nil {
		return nil, nil, err
	}

	if resultMidtrans == nil {
		fmt.Println("Result midtrans : ", resultMidtrans)
		return nil, nil, errors.New("midtrans error, no response payload")
	}

	_, err = bs.d.UpdateBookingStatus(bookingID, booking.Booking{
		BookingStatus: 2,
		VAAccount:     resultMidtrans.RedirectURL,
	})

	header, body := bs.em.HTMLBodyBookingConfirmation(*bookDetail)

	err = bs.em.SendEmail(cust.User.Email, header, body)

	if err != nil {
		return nil, nil, err
	}

	return res, resultMidtrans, nil
}

func (bs *BookingService) GetAllBookingDemo(search string, page uint, pageSize uint, status uint) ([]booking.BookingDemo, uint, uint, error) {
	result, totalPage, totalItems, err := bs.d.GetAllBookingDemo(search, page, pageSize, status)
	if err != nil {
		return nil, 0, 0, errors.New("Get All Booking Demo Process Failed")
	}
	return result, totalPage, totalItems, nil
}

func (bs *BookingService) GetAllBooking(search string, page uint, pageSize uint, status uint) ([]booking.Booking, uint, uint, error) {
	result, totalPage, totalItems, err := bs.d.GetAllBooking(search, page, pageSize, status)
	if err != nil {
		return nil, 0, 0, errors.New("Get All Booking Process Failed")
	}
	return result, totalPage, totalItems, nil
}

func (bs *BookingService) GetBookingByID(id uuid.UUID) (*booking.Booking, error) {
	result, err := bs.d.GetBookingByID(id)
	if err != nil {
		return nil, errors.New("Get Booking By ID Failed")
	}
	return result, nil
}

func (bs *BookingService) GetBookingByBookingInvoice(bookingInvoice string) (*booking.Booking, error) {
	result, err := bs.d.GetBookingByBookingInvoice(bookingInvoice)
	if err != nil {
		return nil, errors.New("Get Booking By Booking Invoice Failed")
	}
	return result, nil
}

func (bs *BookingService) UpdateBooking(id uuid.UUID, newData booking.Booking) (*booking.Booking, *string, error) {
	result, name, err := bs.d.UpdateBooking(id, newData)
	if err != nil {
		return nil, nil, errors.New("Update Booking By ID Failed")
	}
	return result, name, nil
}

func (bs *BookingService) DeleteBooking(id uuid.UUID) (bool, error) {
	result, err := bs.d.DeleteBooking(id)
	if err != nil {
		return false, errors.New("Delete Booking By ID Failed")
	}
	return result, nil
}

func (bs *BookingService) ReportBooking() ([]booking.Booking, error) {
	result, err := bs.d.ReportBooking()
	if err != nil {
		return nil, errors.New("Report booking failed")
	}
	return result, nil
}

func (bs *BookingService) CheckBookingAvailability(officeID uint, startTime time.Time, endTime time.Time) (bool, error) {
	_, err := bs.d.CheckBookingAvailability(officeID, startTime, endTime)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (bs *BookingService) CancelBooking(bookingID uuid.UUID, explanation string) (bool, error) {
	_, err := bs.d.CancelBooking(bookingID, explanation)

	if err != nil {
		return false, err
	}

	return true, nil
}

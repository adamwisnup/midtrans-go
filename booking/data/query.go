package data

import (
	booking "ceo-suite-go/features/booking"
	"ceo-suite-go/features/customer"
	"ceo-suite-go/features/office"
	"ceo-suite-go/features/users"
	"ceo-suite-go/helper"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/example"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	MODULE = "[BOOKING SERVICE]"
)

type BookingData struct {
	db *gorm.DB
	cd customer.CustomerDataInterface
	od office.OfficeDataInterface
}

func New(db *gorm.DB, cd customer.CustomerDataInterface, od office.OfficeDataInterface) booking.BookingDataInterface {
	return &BookingData{
		db: db,
		cd: cd,
		od: od,
	}
}

func (bd *BookingData) GetAllPaymentList() ([]booking.PaymentList, error) {
	var paymentList []booking.PaymentList

	if err := bd.db.Model(&PaymentList{}).Where("deleted_at IS NULL").Find(&paymentList).Error; err != nil {
		return nil, err
	}

	return paymentList, nil
}

func (bd *BookingData) GetPaymentDataByID(id string) (*booking.Payment, error) {
	var payment booking.Payment
	// var err error

	if err := bd.db.Model(&Payment{}).
		Where("payment_invoice = ?", id).
		First(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

func (bd *BookingData) GetPaymentByID(id uint) (*booking.PaymentList, error) {
	var payment booking.PaymentList

	if err := bd.db.Model(&PaymentList{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").Find(&payment).Error; err != nil {
		return nil, err
	}

	return &payment, nil
}

func (bd *BookingData) Checkout(bookingID uuid.UUID) (*booking.Payment, error) {

	bookingData, _ := bd.GetBookingByID(bookingID)
	uuid := uuid.Must(uuid.NewRandom())
	var dbDataPayment = new(Payment)
	dbDataPayment.ID = uuid
	dbDataPayment.BookingID = bookingData.ID
	dbDataPayment.PaymentType = bookingData.PaymentType
	dbDataPayment.PaymentStatus = 5
	dbDataPayment.TotalPrice = bookingData.FinalPrice
	dbDataPayment.PaymentInvoice = "INV-" + example.Random()

	if err := bd.db.Create(dbDataPayment).Error; err != nil {
		return nil, err
	}

	returnData := booking.Payment{
		BookingID:      dbDataPayment.BookingID,
		PaymentStatus:  dbDataPayment.PaymentStatus,
		PaymentType:    dbDataPayment.PaymentType,
		TotalPrice:     dbDataPayment.TotalPrice,
		PaymentInvoice: dbDataPayment.PaymentInvoice,
	}

	return &returnData, nil
}

func (bd *BookingData) BookingSequence() (*int, error) {
	var nextSeq int

	err := bd.db.Raw("SELECT nextval('invoice_seq')").Scan(&nextSeq).Error
	if err != nil {
		return nil, err
	}

	return &nextSeq, nil
}

func (bd *BookingData) CalculateDuration(start, end time.Time) (normalDuration, overtimeDuration time.Duration, err error) {

	if start.Weekday() == time.Saturday {
		return 0, 0, errors.New("you cannot book meeting rooms on the weekend")
	}

	if start.Weekday() == time.Sunday {
		return 0, 0, errors.New("you cannot book meeting rooms on the weekend")
	}

	workingStart := time.Date(start.Year(), start.Month(), start.Day(), 8, 30, 0, 0, start.Location())
	workingEnd := time.Date(start.Year(), start.Month(), start.Day(), 17, 30, 0, 0, start.Location())

	if start.Day() != end.Day() {
		return 0, 0, errors.New("booking a meeting must be in the same date")
	}

	if start.After(workingStart) && end.Before(workingEnd) {
		normalDuration = end.Sub(start)
		return normalDuration, 0, nil
	}

	if start.Before(workingStart) && end.Before(workingStart) {
		overtimeDuration = end.Sub(start)
		return 0, overtimeDuration, nil
	}

	if start.After(workingEnd) && end.After(workingEnd) {
		overtimeDuration = end.Sub(start)
		return 0, overtimeDuration, nil
	}

	var countAwal time.Duration
	var countAkhir time.Duration

	if start.Before(workingStart) {
		countAwal = workingStart.Sub(start)
		fmt.Println("1 Count = ", countAwal)
	}

	if end.After(workingEnd) {
		countAkhir = end.Sub(workingEnd)
		fmt.Println("2 Count = ", countAkhir)
	}

	allDuration := end.Sub(start)
	overtimeDuration = countAwal + countAkhir
	normalDuration = allDuration - overtimeDuration

	return normalDuration, overtimeDuration, nil
}

func (bd *BookingData) GetOfficeName(officeRegion int) (*string, error) {

	var data string

	switch officeRegion {
	case 1:
		data = "SSC"
	case 2:
		data = "OPP"
	case 3:
		data = "GKBI"
	case 4:
		data = "AXA"
	case 5:
		data = "IDX"
	default:
		data = ""
	}

	return &data, nil
}

func (bd *BookingData) CreateBookingGuest(newData booking.Booking, customerID uint) (*booking.Booking, error) {

	uuid := uuid.Must(uuid.NewRandom())
	_, err := bd.CheckBookingAvailability(newData.OfficeID, newData.BookingStartDate, newData.BookingEndDate)
	if err != nil {
		return nil, err
	}

	office, _ := bd.od.GetOfficeByID(newData.OfficeID)

	checkDuration := newData.BookingEndDate.Sub(newData.BookingStartDate)

	if office.OfficeCategory == 2 && checkDuration.Hours() > 23 {
		return nil, errors.New("you cant book a meeting room more than 24 hour")
	}

	if checkDuration.Hours() < 24 {

		duration, overtimeDuration, err := bd.CalculateDuration(newData.BookingStartDate, newData.BookingEndDate)

		if err != nil {
			return nil, err
		}

		fmt.Printf("Duration biasa: %v hour\n", duration.Hours())
		fmt.Printf("Duration overtime: %v hour\n", overtimeDuration.Hours())

		fmt.Printf("Duration biasa ceiled: %v hour\n", math.Ceil(duration.Hours()))
		fmt.Printf("Duration overtime ceiled: %v hour\n", math.Ceil(overtimeDuration.Hours()))

		var pricePerDays, pricePerHour int
		var pricePerHalfDay, pricePerWeek, pricePerMonth int = 0, 0, 0

		if duration.Hours() == 9 {
			pricePerDays = int(office.Price.PriceDaily)
			pricePerHour = 0
		} else {
			pricePerDays = 0
			pricePerHour = int(math.Ceil(duration.Hours())) * int(office.Price.PricePerHour)
		}

		priceOvertime := math.Ceil(overtimeDuration.Hours()) * float64(office.Price.PriceOvertime)

		logrus.Info("Price Overtime: ", priceOvertime)
		logrus.Info("Price per hour count: ", pricePerHour)
		logrus.Info("Price per half day count: ", pricePerHalfDay)
		logrus.Info("Price per days count: ", pricePerDays)
		logrus.Info("Price per week count: ", pricePerWeek)
		logrus.Info("Price per month count: ", pricePerMonth)

		totalPriceWithoutCateringOvertime := uint(pricePerMonth + pricePerWeek + pricePerDays + pricePerHalfDay + pricePerHour)
		totalPrice := totalPriceWithoutCateringOvertime + uint(priceOvertime) + newData.CateringPrice
		ppn := (11 * totalPrice) / 100
		finalPrice := (ppn + totalPrice)

		logrus.Info("Booking Price Minimum: ", office.Price.PriceMinimum)

		if totalPrice < office.Price.PriceMinimum {
			return nil, errors.New("you havent match the minimum price")
		}

		if office.Price.DirectDiscountOffice != nil {
			if office.Price.DirectDiscountOffice.Quantity > 0 {
				if office.Price.DirectDiscountOffice.DiscountDetail.Type == "direct_discount" {
					fmt.Printf("You get a %v discount!", office.Price.DirectDiscountOffice.DiscountDetail.Amount)
					newData.Discount = uint(office.Price.DirectDiscountOffice.DiscountDetail.Amount)
					newData.Notes = fmt.Sprintf("Achieved %v promotion, got a %v%% discount!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				}
			}
		}

		if office.Price.PricePerHour == 0 {
			return nil, errors.New("the duration must be more than 24 hours")
		}

		logrus.Info("Final Price: ", finalPrice)

		seqNum, _ := bd.BookingSequence()
		officeName, _ := bd.GetOfficeName(int(office.RegionID))

		var dbData = new(Booking)
		dbData.ID = uuid
		dbData.BookingInvoice = fmt.Sprintf("%v/%v/%02d/%d", *officeName, *seqNum, time.Now().Month(), time.Now().Year())
		dbData.OfficeID = newData.OfficeID
		dbData.UserID = newData.UserID
		dbData.CustomerID = customerID
		dbData.HourCount = uint(math.Ceil(duration.Hours()))
		if duration.Hours() == 9 {
			dbData.DaysCount = 1
			dbData.HourCount = 0
		}
		dbData.DaysCount = 0
		dbData.WeekCount = 0
		dbData.MonthCount = 0
		dbData.PaymentType = newData.PaymentType
		dbData.BookingPrice = totalPriceWithoutCateringOvertime
		dbData.OvertimeCount = uint(math.Ceil(overtimeDuration.Hours()))
		dbData.OvertimePrice = uint(priceOvertime)
		dbData.CateringPrice = newData.CateringPrice
		dbData.TotalPrice = totalPrice
		dbData.PPN = ppn
		dbData.FinalPrice = finalPrice
		dbData.BookingStatus = BookingCreated
		dbData.BookingStartDate = newData.BookingStartDate
		dbData.BookingEndDate = newData.BookingEndDate
		dbData.BookingExpirationTime = time.Now().Add(10 * time.Minute)
		dbData.UsesDescription = newData.UsesDescription
		dbData.Notes = newData.Notes

		fmt.Println("Dbdata: ", dbData)

		if err := bd.db.Create(dbData).Error; err != nil {
			return nil, err
		}

		resDB := booking.Booking{
			ID: dbData.ID,
		}

		return &resDB, nil
	}

	duration := checkDuration

	daysCount := uint(duration.Hours()) / 24
	remainingDays := daysCount % 30
	remainingHours := uint(duration.Hours()) % 24
	weekCount := remainingDays / 7
	monthCount := daysCount / 30
	halfDayCount := remainingHours / 12
	hourCount := remainingHours - (halfDayCount * 12)
	daysLeft := daysCount - (monthCount * 30) - (weekCount * 7)

	if weekCount == 4 {
		weekCount = 0
		monthCount += 1
	}

	logrus.Info("Price Per Hour Count: ", hourCount)
	logrus.Info("Price Half Day Count: ", halfDayCount)
	logrus.Info("Price Per Day Count: ", daysLeft)
	logrus.Info("Price Per Week Count: ", weekCount)
	logrus.Info("Price Per Month  count: ", monthCount)

	pricePerHour := hourCount * office.Price.PricePerHour
	pricePerDays := daysLeft * office.Price.PriceDaily
	pricePerWeek := weekCount * office.Price.PriceWeekly
	pricePerMonth := monthCount * office.Price.PriceMonthly
	priceOvertime := 0

	logrus.Info("Price Overtime: ", priceOvertime)
	logrus.Info("Price per hour count: ", pricePerHour)
	logrus.Info("Price per days count: ", pricePerDays)
	logrus.Info("Price per week count: ", pricePerWeek)
	logrus.Info("Price per month count: ", pricePerMonth)

	totalPrice := pricePerMonth + pricePerWeek + pricePerDays + pricePerHour

	var finalPrice, deposit uint = 0, 0
	if monthCount > 0 {
		deposit = office.Price.PriceMonthly * 2
		finalPrice = (((11 * totalPrice) / 100) + totalPrice) + deposit
	} else {
		finalPrice = (((11 * totalPrice) / 100) + totalPrice)
	}

	logrus.Info("Booking Price Minimum: ", office.Price.PriceMinimum)

	if totalPrice < office.Price.PriceMinimum {
		return nil, errors.New("you havent match the minimum price")
	}

	if office.Price.DirectDiscountOffice != nil {
		if office.Price.DirectDiscountOffice.Quantity > 0 {
			if office.Price.DirectDiscountOffice.DiscountDetail.Type == "buy_x_get_x" {
				if monthCount >= office.Price.DirectDiscountOffice.BuyX {
					fmt.Printf("You get free %v month!", office.Price.DirectDiscountOffice.GetX)
					monthCount = monthCount + office.Price.DirectDiscountOffice.GetX
					newData.BookingEndDate = newData.BookingEndDate.AddDate(0, int(office.Price.DirectDiscountOffice.GetX), 0)
					newData.Notes = fmt.Sprintf("Achieved %v promotion, got free %v month(s)!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.GetX)
				}
			} else if office.Price.DirectDiscountOffice.DiscountDetail.Type == "direct_discount" {
				fmt.Printf("You get a %v discount!", office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				newData.Discount = uint(office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				newData.Notes = fmt.Sprintf("Achieved %v promotion, got a %v%% discount!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.DiscountDetail.Amount)
			}
		}
	}

	if office.Price.PricePerHour == 0 && duration.Hours() < 24 {
		return nil, errors.New("the duration must be more than 24 hours")
	}
	logrus.Info("Final Price: ", finalPrice)

	seqNum, _ := bd.BookingSequence()
	officeName, _ := bd.GetOfficeName(int(office.RegionID))

	var dbData = new(Booking)
	dbData.ID = uuid
	dbData.BookingInvoice = fmt.Sprintf("%v/%v/%02d/%d", *officeName, *seqNum, time.Now().Month(), time.Now().Year())
	dbData.OfficeID = newData.OfficeID
	dbData.UserID = newData.UserID
	dbData.CustomerID = customerID
	dbData.BookingStartDate = newData.BookingStartDate
	dbData.BookingEndDate = newData.BookingEndDate
	dbData.BookingExpirationTime = time.Now().Add(15 * time.Minute)
	dbData.MonthCount = monthCount
	dbData.WeekCount = weekCount
	dbData.DaysCount = daysLeft
	dbData.HourCount = hourCount
	dbData.PaymentType = newData.PaymentType
	dbData.BookingPrice = newData.BookingPrice
	dbData.OvertimeCount = 0
	dbData.OvertimePrice = 0
	dbData.CateringPrice = 0
	dbData.BookingStatus = BookingCreated
	dbData.UsesDescription = newData.UsesDescription
	dbData.Notes = newData.Notes
	dbData.BookingPrice = totalPrice
	dbData.TotalPrice = totalPrice
	dbData.PPN = (11 * totalPrice) / 100
	dbData.Deposit = deposit
	dbData.FinalPrice = finalPrice

	if err := bd.db.Create(dbData).Error; err != nil {
		return nil, err
	}

	resDB := booking.Booking{
		ID: dbData.ID,
	}

	return &resDB, nil
}

func (bd *BookingData) CheckBookingAvailability(officeID uint, startTime time.Time, endTime time.Time) (bool, error) {
	var countBooking int64

	office, _ := bd.od.GetOfficeByID(officeID)

	checkDuration := endTime.Sub(startTime)

	if office.OfficeCategory == 2 && checkDuration.Hours() > 23 {
		return false, errors.New("you cant book a meeting room more than 24 hour")
	}

	if office.Price.PricePerHour == 0 && checkDuration.Hours() < 24 {
		return false, errors.New("the duration must be more than 24 hours")
	}

	if startTime.Equal(endTime) {
		return false, errors.New("start time and end time cannot be the same")
	}

	if startTime.Before(time.Now()) {
		return false, errors.New("you cannot book in past time")
	}

	if endTime.Before(startTime) {
		return false, errors.New("your end time book must be greater than start time")
	}

	err := bd.db.Model(&Booking{}).
		Where("office_id = ?", officeID).
		Where("(booking_start_date < ? AND booking_end_date > ?) OR (booking_start_date < ? AND booking_end_date > ?)", endTime, startTime, endTime, startTime).
		Where("booking_expiration_time > ?", time.Now()).
		Where("deleted_at IS NULL").
		Count(&countBooking).
		Error

	if err != nil {
		return false, err
	}

	fmt.Printf("COUNT %v \n", countBooking)

	if countBooking > 0 {
		return false, errors.New("there is already a booking in that date")
	}

	return true, nil
}

func (bd *BookingData) CreateBooking(newData booking.Booking) (*booking.Booking, *booking.Payment, *string, error) {

	uuid := uuid.Must(uuid.NewRandom())
	_, err := bd.CheckBookingAvailability(newData.OfficeID, newData.BookingStartDate, newData.BookingEndDate)
	if err != nil {
		return nil, nil, nil, err
	}

	office, _ := bd.od.GetOfficeByID(newData.OfficeID)
	custID, err := bd.cd.GetCustomerByUserID(newData.UserID)
	if err != nil {
		return nil, nil, nil, err
	}

	checkDuration := newData.BookingEndDate.Sub(newData.BookingStartDate)

	if office.OfficeCategory == 2 && checkDuration.Hours() > 23 {
		return nil, nil, nil, errors.New("you cant book a meeting room more than 24 hour")
	}

	if checkDuration.Hours() < 24 {

		fmt.Println("Kena disini")

		duration, overtimeDuration, err := bd.CalculateDuration(newData.BookingStartDate, newData.BookingEndDate)

		if err != nil {
			return nil, nil, nil, err
		}

		fmt.Printf("Duration biasa: %v hour\n", duration.Hours())
		fmt.Printf("Duration overtime: %v hour\n", overtimeDuration.Hours())

		fmt.Printf("Duration biasa ceiled: %v hour\n", math.Ceil(duration.Hours()))
		fmt.Printf("Duration overtime ceiled: %v hour\n", math.Ceil(overtimeDuration.Hours()))

		var pricePerDays, pricePerHour int
		var pricePerHalfDay, pricePerWeek, pricePerMonth int = 0, 0, 0

		if duration.Hours() == 9 {
			pricePerDays = int(office.Price.PriceDaily)
			pricePerHour = 0
		} else {
			pricePerDays = 0
			pricePerHour = int(math.Ceil(duration.Hours())) * int(office.Price.PricePerHour)
		}

		priceOvertime := math.Ceil(overtimeDuration.Hours()) * float64(office.Price.PriceOvertime)

		logrus.Info("Price Overtime: ", priceOvertime)
		logrus.Info("Price per hour count: ", pricePerHour)
		logrus.Info("Price per half day count: ", pricePerHalfDay)
		logrus.Info("Price per days count: ", pricePerDays)
		logrus.Info("Price per week count: ", pricePerWeek)
		logrus.Info("Price per month count: ", pricePerMonth)

		totalPriceWithoutCateringOvertime := uint(pricePerMonth + pricePerWeek + pricePerDays + pricePerHalfDay + pricePerHour)
		totalPrice := totalPriceWithoutCateringOvertime + uint(priceOvertime) + newData.CateringPrice
		ppn := (11 * totalPrice) / 100
		finalPrice := (ppn + totalPrice)

		logrus.Info("Booking Price Minimum: ", office.Price.PriceMinimum)

		if totalPrice < office.Price.PriceMinimum {
			return nil, nil, nil, errors.New("you havent match the minimum price")
		}

		if office.Price.DirectDiscountOffice != nil {
			if office.Price.DirectDiscountOffice.Quantity > 0 {
				if office.Price.DirectDiscountOffice.DiscountDetail.Type == "direct_discount" {
					fmt.Printf("You get a %v discount!", office.Price.DirectDiscountOffice.DiscountDetail.Amount)
					newData.Discount = uint(office.Price.DirectDiscountOffice.DiscountDetail.Amount)
					newData.Notes = fmt.Sprintf("Achieved %v promotion, got a %v%% discount!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				}
			}
		}

		if office.Price.PricePerHour == 0 {
			return nil, nil, nil, errors.New("the duration must be more than 24 hours")
		}

		logrus.Info("Final Price: ", finalPrice)

		seqNum, _ := bd.BookingSequence()
		officeName, _ := bd.GetOfficeName(int(office.RegionID))

		var dbData = new(Booking)
		dbData.ID = uuid
		dbData.BookingInvoice = fmt.Sprintf("%v/%v/%02d/%d", *officeName, *seqNum, time.Now().Month(), time.Now().Year())
		dbData.OfficeID = newData.OfficeID
		dbData.UserID = newData.UserID
		dbData.CustomerID = custID.ID
		dbData.HourCount = uint(math.Ceil(duration.Hours()))
		if duration.Hours() == 9 {
			dbData.DaysCount = 1
			dbData.HourCount = 0
		}
		dbData.DaysCount = 0
		dbData.WeekCount = 0
		dbData.MonthCount = 0
		dbData.PaymentType = newData.PaymentType
		dbData.BookingPrice = totalPriceWithoutCateringOvertime
		dbData.OvertimeCount = uint(math.Ceil(overtimeDuration.Hours()))
		dbData.OvertimePrice = uint(priceOvertime)
		dbData.CateringPrice = newData.CateringPrice
		dbData.TotalPrice = totalPrice
		dbData.PPN = ppn
		dbData.FinalPrice = finalPrice
		dbData.BookingStatus = BookingCreated
		dbData.BookingStartDate = newData.BookingStartDate
		dbData.BookingEndDate = newData.BookingEndDate
		dbData.BookingExpirationTime = time.Now().Add(10 * time.Minute)
		dbData.UsesDescription = newData.UsesDescription
		dbData.Notes = newData.Notes

		fmt.Println("Dbdata: ", dbData)

		if err := bd.db.Create(dbData).Error; err != nil {
			return nil, nil, nil, err
		}

		resDB := booking.Booking{
			ID: dbData.ID,
		}

		return &resDB, nil, &custID.FullName, nil
	}

	duration := checkDuration

	daysCount := uint(duration.Hours()) / 24
	remainingDays := daysCount % 30
	remainingHours := uint(duration.Hours()) % 24
	weekCount := remainingDays / 7
	monthCount := daysCount / 30
	halfDayCount := remainingHours / 12
	hourCount := remainingHours - (halfDayCount * 12)
	daysLeft := daysCount - (monthCount * 30) - (weekCount * 7)

	if weekCount == 4 {
		weekCount = 0
		monthCount += 1
	}

	logrus.Info("Price Per Hour Count: ", hourCount)
	logrus.Info("Price Half Day Count: ", halfDayCount)
	logrus.Info("Price Per Day Count: ", daysLeft)
	logrus.Info("Price Per Week Count: ", weekCount)
	logrus.Info("Price Per Month  count: ", monthCount)

	pricePerHour := hourCount * office.Price.PricePerHour
	pricePerDays := daysLeft * office.Price.PriceDaily
	pricePerWeek := weekCount * office.Price.PriceWeekly
	pricePerMonth := monthCount * office.Price.PriceMonthly
	priceOvertime := 0

	logrus.Info("Price Overtime: ", priceOvertime)
	logrus.Info("Price per hour count: ", pricePerHour)
	logrus.Info("Price per days count: ", pricePerDays)
	logrus.Info("Price per week count: ", pricePerWeek)
	logrus.Info("Price per month count: ", pricePerMonth)

	totalPrice := pricePerMonth + pricePerWeek + pricePerDays + pricePerHour

	var finalPrice, deposit uint = 0, 0
	if monthCount > 0 {
		deposit = office.Price.PriceMonthly * 2
		finalPrice = (((11 * totalPrice) / 100) + totalPrice) + deposit
	} else {
		finalPrice = (((11 * totalPrice) / 100) + totalPrice)
	}

	logrus.Info("Booking Price Minimum: ", office.Price.PriceMinimum)

	if totalPrice < office.Price.PriceMinimum {
		return nil, nil, nil, errors.New("you havent match the minimum price")
	}

	if office.Price.DirectDiscountOffice != nil {
		if office.Price.DirectDiscountOffice.Quantity > 0 {
			if office.Price.DirectDiscountOffice.DiscountDetail.Type == "buy_x_get_x" {
				if monthCount >= office.Price.DirectDiscountOffice.BuyX {
					fmt.Printf("You get free %v month!", office.Price.DirectDiscountOffice.GetX)
					monthCount = monthCount + office.Price.DirectDiscountOffice.GetX
					newData.BookingEndDate = newData.BookingEndDate.AddDate(0, int(office.Price.DirectDiscountOffice.GetX), 0)
					newData.Notes = fmt.Sprintf("Achieved %v promotion, got free %v month(s)!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.GetX)
				}
			} else if office.Price.DirectDiscountOffice.DiscountDetail.Type == "direct_discount" {
				fmt.Printf("You get a %v discount!", office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				newData.Discount = uint(office.Price.DirectDiscountOffice.DiscountDetail.Amount)
				newData.Notes = fmt.Sprintf("Achieved %v promotion, got a %v%% discount!", office.Price.DirectDiscountOffice.DiscountDetail.Name, office.Price.DirectDiscountOffice.DiscountDetail.Amount)
			}
		}
	}

	if office.Price.PricePerHour == 0 && duration.Hours() < 24 {
		return nil, nil, nil, errors.New("the duration must be more than 24 hours")
	}
	logrus.Info("Final Price: ", finalPrice)

	seqNum, _ := bd.BookingSequence()
	officeName, _ := bd.GetOfficeName(int(office.RegionID))

	var dbData = new(Booking)
	dbData.ID = uuid
	dbData.BookingInvoice = fmt.Sprintf("%v/%v/%02d/%d", *officeName, *seqNum, time.Now().Month(), time.Now().Year())
	dbData.OfficeID = newData.OfficeID
	dbData.UserID = newData.UserID
	dbData.CustomerID = custID.ID
	dbData.BookingStartDate = newData.BookingStartDate
	dbData.BookingEndDate = newData.BookingEndDate
	dbData.BookingExpirationTime = time.Now().Add(15 * time.Minute)
	dbData.MonthCount = monthCount
	dbData.WeekCount = weekCount
	dbData.DaysCount = daysLeft
	dbData.HourCount = hourCount
	dbData.PaymentType = newData.PaymentType
	dbData.BookingPrice = newData.BookingPrice
	dbData.OvertimeCount = 0
	dbData.OvertimePrice = 0
	dbData.CateringPrice = 0
	dbData.BookingStatus = BookingCreated
	dbData.UsesDescription = newData.UsesDescription
	dbData.Notes = newData.Notes
	dbData.BookingPrice = totalPrice
	dbData.TotalPrice = totalPrice
	dbData.PPN = (11 * totalPrice) / 100
	dbData.Deposit = deposit

	dbData.FinalPrice = finalPrice

	if err := bd.db.Create(dbData).Error; err != nil {
		return nil, nil, nil, err
	}

	resDB := booking.Booking{
		ID: dbData.ID,
	}

	return &resDB, nil, &custID.FullName, nil
}

func (bd *BookingData) CreateCateringDetails(newData booking.CateringDetails) (*booking.CateringDetails, error) {
	var dbData = new(CateringDetails)

	dbData.BookingID = newData.BookingID
	dbData.CateringID = newData.CateringID
	dbData.CateringCountPax = newData.CateringCountPax
	dbData.CateringPrice = newData.CateringPrice

	if err := bd.db.Create(dbData).Error; err != nil {
		return nil, err
	}

	resDB := booking.CateringDetails{
		BookingID:        dbData.BookingID,
		CateringID:       dbData.CateringID,
		CateringCountPax: dbData.CateringCountPax,
		CateringPrice:    dbData.CateringPrice,
	}

	return &resDB, nil
}

func (bd *BookingData) GetAllBookingDemo(search string, page uint, pageSize uint, status uint) ([]booking.BookingDemo, uint, uint, error) {
	var booking []booking.BookingDemo

	qry := bd.db.Table("bookings").
		Select("bookings.*, users.name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id")

	if status > 0 {
		qry = qry.Where("bookings.booking_status = ?", status)
	}

	if search != "" {
		qry = qry.Where("bookings.booking_invoice ILIKE ?", search)
	}

	qry = qry.Where("bookings.deleted_at IS NULL").
		Order("bookings.created_at DESC")

	paginatedQuery, totalPage, totalItems, err := helper.PaginateQuery(qry, page, pageSize)

	if err != nil {
		return nil, 0, 0, err
	}

	err = paginatedQuery.Find(&booking).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return booking, totalPage, totalItems, err
}

func (bd *BookingData) GetAllBooking(search string, page uint, pageSize uint, status uint) ([]booking.Booking, uint, uint, error) {
	var booking []booking.Booking

	baseQuery := bd.db.Model(&Booking{}).
		Select("bookings.*, users.name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id")

	if status > 0 {
		baseQuery = baseQuery.Where("bookings.booking_status = ?", status)
	}

	if search != "" {
		baseQuery = baseQuery.Where("bookings.booking_invoice ILIKE ?", search)
	}

	baseQuery = baseQuery.Where("bookings.deleted_at IS NULL").
		Order("bookings.created_at DESC").
		Preload("Payment")

	paginatedQuery, totalPage, totalItems, err := helper.PaginateQuery(baseQuery, page, pageSize)

	if err != nil {
		return nil, 0, 0, err
	}

	err = paginatedQuery.Find(&booking).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return booking, totalPage, totalItems, err
}

func (bd *BookingData) GetBookingByID(id uuid.UUID) (*booking.Booking, error) {
	var booking booking.Booking

	err := bd.db.Model(&Booking{}).
		Select("bookings.*, users.name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id").
		Where("bookings.id = ?", id).
		Where("bookings.deleted_at IS NULL").
		Preload("Payment").
		Find(&booking).
		Error

	return &booking, err
}

func (bd *BookingData) GetBookingByBookingInvoice(bookingInvoice string) (*booking.Booking, error) {
	var booking booking.Booking

	err := bd.db.Model(&Booking{}).
		Select("bookings.*, users.name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id").
		Where("bookings.booking_invoice = ?", bookingInvoice).
		Where("bookings.deleted_at IS NULL").
		Preload("Payment").
		Find(&booking).
		Error

	return &booking, err
}

func (bd *BookingData) GetCustomersBooking(userID uint, status uint) ([]booking.BookingCustomer, error) {
	var booking []booking.BookingCustomer

	qry := bd.db.Model(&Booking{}).
		Select("bookings.*, users.name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id").
		// Preload("Payment").
		Where("bookings.user_id = ?", userID)

	if status > 0 {
		qry = qry.Where("bookings.booking_status = ?", status)
	}

	err := qry.Where("bookings.deleted_at IS NULL").
		Order("bookings.created_at DESC").
		Find(&booking).Error

	if err != nil {
		return nil, err
	}

	return booking, err
}

func (bd *BookingData) GetCustomersBookingByID(id uuid.UUID) (*booking.BookingCustomerDetails, error) {
	var book booking.BookingCustomerDetails

	err := bd.db.Model(&Booking{}).
		Select("bookings.*, customers.full_name as customer_name, offices.name as office_name").
		Joins("LEFT JOIN customers ON customers.id = bookings.customer_id").
		Joins("LEFT JOIN offices ON offices.id = bookings.office_id").
		Where("bookings.id = ?", id).
		Where("bookings.deleted_at IS NULL").
		Preload("CateringDetails").
		Preload("Payment").
		Find(&book).
		Error

	resOffice, _ := bd.od.GetOfficeByID(book.OfficeID)

	if book.ID == uuid.Nil {
		return nil, errors.New("failed to get booking data")
	}

	if time.Now().After(book.BookingExpirationTime) && book.BookingStatus < 3 {
		_, err = bd.UpdateBookingStatus(id, booking.Booking{
			BookingStatus: BookingExpired,
		})

		book.BookingStatus = BookingExpired
	}

	for i := range book.CateringDetails {
		cateringData, _ := bd.od.GetCateringByID(book.CateringDetails[i].CateringID)

		book.CateringDetails[i].CateringName = cateringData.Name
		fmt.Println("Catering data: ", cateringData)
		cateringUnitPrice := book.CateringDetails[i].CateringPrice / book.CateringDetails[i].CateringCountPax
		// fmt.Println("Catering data: ", catering)
		book.CateringDetails[i].CateringUnitPrice = &cateringUnitPrice
	}

	if book.Payment != nil {
		book.Payment.ExpiredAt = book.BookingExpirationTime
	}

	if book.MonthCount > 0 {
		book.Deposit = (resOffice.Price.PriceMonthly * 2)
	}

	return &book, err
}

func (bd *BookingData) GetAndUpdatePayment(newData booking.Payment, id string) (bool, error) {
	var payment booking.Payment
	var err error

	if err = bd.db.Model(&Payment{}).
		Where("payment_invoice = ?", id).
		First(&payment).Error; err != nil {
		return false, err
	}

	_, err = bd.UpdatePaymentByBookingID(payment.BookingID, booking.Payment{
		PaymentStatus: newData.PaymentStatus,
	})

	if err != nil {
		return false, err
	}

	if newData.PaymentStatus == 2 {
		_, err = bd.UpdateBookingStatus(payment.BookingID, booking.Booking{
			BookingStatus: 3,
		})
	} else if newData.PaymentStatus == 3 {
		_, err = bd.UpdateBookingStatus(payment.BookingID, booking.Booking{
			BookingStatus: 4,
		})
	} else if newData.PaymentStatus == 4 {
		_, err = bd.UpdateBookingStatus(payment.BookingID, booking.Booking{
			BookingStatus: 4,
		})
	} else if newData.PaymentStatus == 5 {
		_, err = bd.UpdateBookingStatus(payment.BookingID, booking.Booking{
			BookingStatus: 2,
		})
	} else {
		_, err = bd.UpdateBookingStatus(payment.BookingID, booking.Booking{
			BookingStatus: 2,
		})
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (bd *BookingData) UpdateBooking(id uuid.UUID, newData booking.Booking) (*booking.Booking, *string, error) {
	var officePrice office.Price
	var existingBooking booking.Booking
	var user users.User
	var office office.Office

	if err := bd.db.First(&existingBooking, id).Error; err != nil {
		return nil, nil, err
	}

	if err := bd.db.First(&office, existingBooking.OfficeID).Error; err != nil {
		return nil, nil, err
	}

	if err := bd.db.First(&officePrice, "office_id = ?", office.ID).Error; err != nil {
		return nil, nil, err
	}

	existingBooking.OfficeID = newData.OfficeID
	existingBooking.BookingStartDate = newData.BookingStartDate
	existingBooking.BookingEndDate = newData.BookingEndDate
	existingBooking.PaymentType = newData.PaymentType
	existingBooking.UsesDescription = newData.UsesDescription

	duration := existingBooking.BookingEndDate.Sub(existingBooking.BookingStartDate)
	daysCount := uint(duration.Hours()) / 24
	remainingHours := uint(duration.Hours()) % 24

	// halfDayCount := remainingHours / 12
	hourCount := remainingHours % 12

	pricePerHour := hourCount * officePrice.PricePerHour
	// pricePerHalfDay := halfDayCount * officePrice.PriceHalfDay
	pricePerDays := daysCount * officePrice.PriceDaily

	logrus.Info("Price per hour count: ", pricePerHour)
	// logrus.Info("Price per half day count: ", pricePerHalfDay)
	logrus.Info("Price per days count: ", pricePerDays)

	totalPrice := pricePerDays + pricePerHour

	if office.Discount > 0 {
		discount := (totalPrice * office.Discount) / 100
		totalPrice -= discount
	}

	existingBooking.TotalPrice = totalPrice
	existingBooking.DaysCount = daysCount
	existingBooking.HourCount = remainingHours
	existingBooking.Discount = office.Discount

	if err := bd.db.Model(&existingBooking).Updates(existingBooking).Error; err != nil {
		return nil, nil, err
	}

	logrus.Info("Total Price: ", existingBooking.TotalPrice)
	logrus.Info("Days Count: ", existingBooking.DaysCount)
	logrus.Info("Hour Count: ", existingBooking.HourCount)
	logrus.Info("Discount: ", existingBooking.Discount)

	if err := bd.db.Model(&user).Where("id = ?", existingBooking.UserID).Select("name").Scan(&user).Error; err != nil {
		return nil, nil, err
	}

	payment := &Payment{
		BookingID:      existingBooking.ID,
		PaymentStatus:  existingBooking.BookingStatus,
		TotalPrice:     existingBooking.TotalPrice,
		PaymentInvoice: "INV-" + example.Random(),
	}

	logrus.Info("Booking ID: ", payment.BookingID)
	logrus.Info("Payment Status: ", payment.PaymentStatus)
	logrus.Info("Total Price: ", payment.TotalPrice)
	logrus.Info("Payment Invoice: ", payment.PaymentInvoice)

	if err := bd.db.Save(payment).Error; err != nil {
		return nil, nil, err
	}

	return &existingBooking, &user.Name, nil
}

func (bd *BookingData) UpdateExpiredBookings() (int64, error) {
	qry := bd.db.Model(&Booking{}).
		Where("booking_status < 3").
		Where("booking_expiration_time < ?", time.Now()).
		Update("booking_status", BookingExpired)

	if qry.Error != nil {
		return 0, qry.Error
	}

	return qry.RowsAffected, nil
}

func (bd *BookingData) UpdateBookingStatus(id uuid.UUID, newData booking.Booking) (bool, error) {

	if err := bd.db.Model(&Booking{}).Where("id = ?", id).Updates(Booking{
		BookingStatus: newData.BookingStatus,
		VAAccount:     newData.VAAccount,
		ReasonCancel:  newData.ReasonCancel,
	}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (bd *BookingData) UpdatePaymentByBookingID(bookingID uuid.UUID, newData booking.Payment) (bool, error) {

	if err := bd.db.Model(&Payment{}).Where("booking_id = ?", bookingID).Updates(Payment{
		PaymentStatus: newData.PaymentStatus,
	}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (bd *BookingData) DeleteBooking(id uuid.UUID) (bool, error) {
	var booking booking.Booking
	var qry *gorm.DB

	qry = bd.db.Table("bookings").
		Select("bookings.*").
		Where("bookings.id = ?", id).
		Where("bookings.deleted_at IS NULL").
		Scan(&booking)

	if err := qry.Error; err != nil {
		return false, err
	}

	if booking.ID == uuid.Nil {
		logrus.Info("Masuk sini, isinya: ", booking)
		return false, errors.New("id not found")
	}

	qry = bd.db.Delete(&Payment{}, "booking_id = ?", id)
	if err := qry.Error; err != nil {
		return false, err
	}

	qry = bd.db.Delete(&Booking{}, "id = ?", id)
	if err := qry.Error; err != nil {
		return false, err
	}

	return true, nil
}

func (bd *BookingData) ReportBooking() ([]booking.Booking, error) {
	var bookings []booking.Booking
	result := bd.db.Find(&bookings)
	return bookings, result.Error
}

func (bd *BookingData) CancelBooking(id uuid.UUID, explanation string) (bool, error) {

	book, err := bd.GetCustomersBookingByID(id)

	if err != nil {
		return false, errors.New("failed to get booking data")
	}

	if book.BookingStatus > 2 {
		return false, errors.New("cannot cancel booking")
	}

	_, err = bd.UpdateBookingStatus(id, booking.Booking{
		BookingStatus: BookingCanceledByCustomer,
		ReasonCancel:  explanation,
	})

	if err != nil {
		return false, errors.New("cannot update booking data")
	}

	return true, nil
}

func (bd *BookingData) CancelBookingByAdmin(id uuid.UUID, explanation string) (bool, error) {

	book, err := bd.GetCustomersBookingByID(id)

	if err != nil {
		return false, errors.New("failed to get booking data")
	}

	if book.BookingStatus > 2 {
		return false, errors.New("cannot cancel booking")
	}

	_, err = bd.UpdateBookingStatus(id, booking.Booking{
		BookingStatus: BookingCanceledByAdmin,
		ReasonCancel:  explanation,
	})

	if err != nil {
		return false, errors.New("cannot update booking data")
	}

	return true, nil
}

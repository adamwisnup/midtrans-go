package handler

import (
	booking "ceo-suite-go/features/booking"
	"ceo-suite-go/features/office"
	"ceo-suite-go/helper"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type BookingHandler struct {
	s   booking.BookingServiceInterface
	os  office.OfficeServiceInterface
	jwt helper.JWTInterface
}

func NewHandler(service booking.BookingServiceInterface, jwt helper.JWTInterface, os office.OfficeServiceInterface) booking.BookingHandlerInterface {
	return &BookingHandler{
		s:   service,
		jwt: jwt,
		os:  os,
	}
}

func (bh *BookingHandler) NotifBooking() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notificationPayload map[string]interface{}

		err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload)

		fmt.Println("Notification Payload:", notificationPayload)

		if err != nil {
			if strings.Contains(err.Error(), "Order ID Not Found") {
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Order ID Not Found", nil))
			}

			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Midtrans POST method error", nil))
		}

		res, err := bh.s.NotifBooking(notificationPayload, booking.Payment{})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Midtrans cannot update the database", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Update", res))
	}
}

func (bh *BookingHandler) GetAllPaymentList() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, err := bh.s.GetAllPaymentList()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success get all payment", res))
	}
}

func (bh *BookingHandler) GetPaymentByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramID = c.Param("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			c.Logger().Error("Handler : Param ID Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid User Input Param ID", nil))
		}

		res, err := bh.s.GetPaymentByID(uint(id))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success get all payment", res))
	}
}

func (bh *BookingHandler) CreateBooking() echo.HandlerFunc {
	return func(c echo.Context) error {

		getID, err := bh.jwt.GetID(c)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "Fail to get id from jwt", nil))
		}

		var input = new(InputRequestTest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		if !input.AcceptTermsConditions {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "you must accept terms and conditions", nil))
		}

		getOffice, err := bh.os.GetOfficeByID(input.OfficeID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		fmt.Println("Before localize: START ", input.BookingStartDate, " END ", input.BookingEndDate)

		localizeStartDate := input.BookingStartDate
		localizeEndDate := input.BookingEndDate

		fmt.Println("After localize: START ", localizeStartDate, " END ", localizeEndDate)

		_, err = bh.s.CheckBookingAvailability(input.OfficeID, localizeStartDate, localizeEndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}
		_, overtimeDuration, err := bh.s.CalculateDuration(localizeStartDate, localizeEndDate)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		if input.CateringRequest != nil && getOffice.OfficeCategory != 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "you only can order catering when booking meeting room", nil))
		}

		var totalPriceCatering uint

		for _, catID := range input.CateringRequest {
			catering, _ := bh.os.GetCateringByID(catID.CateringID)
			if catering.Category == 10 {
				totalPriceCatering = totalPriceCatering + catering.CateringPrice.PricePerEvent
			} else {
				priceEachCat := catering.CateringPrice.PricePerPax * catID.CateringPaxCount
				totalPriceCatering = totalPriceCatering + priceEachCat
			}
		}

		if input.UseOvertimeAC && uint(math.Ceil((overtimeDuration.Hours()))) > 0 {
			if getOffice.RegionID != 2 {
				addOn, _ := bh.os.GetCateringByID(8)
				totalPriceCatering = totalPriceCatering + uint(math.Ceil((overtimeDuration.Hours())))*addOn.CateringPrice.PricePerPax
				addOnDetails := CateringRequest{
					CateringID:       8,
					CateringPaxCount: uint(math.Ceil((overtimeDuration.Hours()))),
				}
				fmt.Println("Total price catering 1001: ", uint(math.Ceil((overtimeDuration.Hours()))), " ", totalPriceCatering)
				input.CateringRequest = append(input.CateringRequest, addOnDetails)
			} else {
				addOn, _ := bh.os.GetCateringByID(9)
				totalPriceCatering = totalPriceCatering + (uint(math.Ceil((overtimeDuration.Hours()))) * addOn.CateringPrice.PricePerPax)
				addOnDetails := CateringRequest{
					CateringID:       9,
					CateringPaxCount: uint(math.Ceil((overtimeDuration.Hours()))),
				}
				fmt.Println("Total price catering 1001: ", uint(math.Ceil((overtimeDuration.Hours()))), " ", totalPriceCatering)
				input.CateringRequest = append(input.CateringRequest, addOnDetails)
			}
		}

		var serviceAddBooking = new(booking.Booking)

		serviceAddBooking.UserID = getID
		serviceAddBooking.OfficeID = input.OfficeID
		serviceAddBooking.CateringPrice = totalPriceCatering
		serviceAddBooking.BookingStartDate = localizeStartDate
		serviceAddBooking.BookingEndDate = localizeEndDate
		serviceAddBooking.UsesDescription = input.UsesDescription

		fmt.Println("Check catering price: ", serviceAddBooking.CateringPrice)

		resultAddedBooking, _, _, err := bh.s.CreateBooking(*serviceAddBooking)

		if err != nil {
			c.Logger().Error("Handler: Input Process Error (CreateBooking): ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		for _, catDetails := range input.CateringRequest {
			var serviceAddCateringDetails = new(booking.CateringDetails)
			// var cateringUnitPrice uint = 0

			catering, _ := bh.os.GetCateringByID(catDetails.CateringID)

			serviceAddCateringDetails.BookingID = resultAddedBooking.ID
			serviceAddCateringDetails.CateringID = catDetails.CateringID
			serviceAddCateringDetails.CateringCountPax = catDetails.CateringPaxCount
			if catering.Category == 10 {
				serviceAddCateringDetails.CateringPrice = catering.CateringPrice.PricePerEvent
			} else {
				serviceAddCateringDetails.CateringPrice = catering.CateringPrice.PricePerPax * catDetails.CateringPaxCount
			}

			_, err := bh.s.CreateCateringDetails(*serviceAddCateringDetails)

			if err != nil {
				c.Logger().Error("Handler: Input Process Error (Create Catering Details): ", err.Error())
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
			}
		}

		var response = new(InputResponse)
		response.BookingID = resultAddedBooking.ID

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success create Booking", response))
	}
}

func (bh *BookingHandler) CreateBookingGuest() echo.HandlerFunc {
	return func(c echo.Context) error {

		var input = new(BookingGuestRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		if !input.AcceptTermsConditions {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "you must accept terms and conditions", nil))
		}

		getOffice, err := bh.os.GetOfficeByID(input.OfficeID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		_, err = bh.s.CheckBookingAvailability(input.OfficeID, input.BookingStartDate, input.BookingEndDate)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		_, overtimeDuration, err := bh.s.CalculateDuration(input.BookingStartDate, input.BookingEndDate)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		if input.CateringRequest != nil && getOffice.OfficeCategory != 2 {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "you only can order catering when booking meeting room", nil))
		}

		var totalPriceCatering uint

		for _, catID := range input.CateringRequest {
			catering, _ := bh.os.GetCateringByID(catID.CateringID)
			if catering.Category == 10 {
				totalPriceCatering = totalPriceCatering + catering.CateringPrice.PricePerEvent
			} else {
				priceEachCat := catering.CateringPrice.PricePerPax * catID.CateringPaxCount
				totalPriceCatering = totalPriceCatering + priceEachCat
			}
		}

		if input.UseOvertimeAC && uint(math.Ceil((overtimeDuration.Hours()))) > 0 {
			if getOffice.RegionID != 2 {
				addOn, _ := bh.os.GetCateringByID(8)
				totalPriceCatering = totalPriceCatering + (uint(math.Ceil((overtimeDuration.Hours()))) * addOn.CateringPrice.PricePerEvent)
				addOnDetails := CateringRequest{
					CateringID:       8,
					CateringPaxCount: uint(math.Ceil((overtimeDuration.Hours()))),
				}
				fmt.Println("Total price catering 1001: ", uint(math.Ceil((overtimeDuration.Hours()))), " ", totalPriceCatering)
				input.CateringRequest = append(input.CateringRequest, addOnDetails)
			} else {
				addOn, _ := bh.os.GetCateringByID(9)
				totalPriceCatering = totalPriceCatering + (uint(math.Ceil((overtimeDuration.Hours()))) * addOn.CateringPrice.PricePerEvent)
				addOnDetails := CateringRequest{
					CateringID:       9,
					CateringPaxCount: uint(math.Ceil((overtimeDuration.Hours()))),
				}
				fmt.Println("Total price catering 1001: ", uint(math.Ceil((overtimeDuration.Hours()))), " ", totalPriceCatering)
				input.CateringRequest = append(input.CateringRequest, addOnDetails)
			}
		}

		var serviceAddGuest = new(booking.BookingGuestRequest)
		serviceAddGuest.FullName = input.FullName
		serviceAddGuest.Email = input.Email
		serviceAddGuest.PhoneNumber = input.PhoneNumber

		var serviceAddBooking = new(booking.Booking)

		serviceAddBooking.OfficeID = input.OfficeID
		serviceAddBooking.CateringPrice = totalPriceCatering
		serviceAddBooking.BookingStartDate = input.BookingStartDate
		serviceAddBooking.BookingEndDate = input.BookingEndDate
		serviceAddBooking.UsesDescription = input.UsesDescription

		resultAddedBooking, customer, err := bh.s.CreateBookingGuest(*serviceAddBooking, *serviceAddGuest)

		if err != nil {
			c.Logger().Error("Handler: Input Process Error (CreateBooking): ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		for _, catDetails := range input.CateringRequest {
			var serviceAddCateringDetails = new(booking.CateringDetails)

			catering, _ := bh.os.GetCateringByID(catDetails.CateringID)

			serviceAddCateringDetails.BookingID = resultAddedBooking.ID
			serviceAddCateringDetails.CateringID = catDetails.CateringID
			serviceAddCateringDetails.CateringCountPax = catDetails.CateringPaxCount
			if catering.Category == 10 {
				serviceAddCateringDetails.CateringPrice = catering.CateringPrice.PricePerEvent
			} else {
				serviceAddCateringDetails.CateringPrice = catering.CateringPrice.PricePerPax * catDetails.CateringPaxCount
			}

			_, err := bh.s.CreateCateringDetails(*serviceAddCateringDetails)

			if err != nil {
				c.Logger().Error("Handler: Input Process Error (Create Catering Details): ", err.Error())
				return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
			}
		}

		var response = new(BookingGuestResponse)
		response.IsGuestMode = true
		response.Customer = customer
		response.BookingID = resultAddedBooking.ID

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success create Booking", response))
	}
}

func (bh *BookingHandler) Checkout() echo.HandlerFunc {
	return func(c echo.Context) error {

		var input = new(CheckoutRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		res, mt, err := bh.s.Checkout(input.BookingID)

		if err != nil {
			c.Logger().Error("Handler: Input Process Error (Checkout): ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		fmt.Println("Res payment", res)
		fmt.Println("Res midtrans", mt)

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success checkout", mt))
	}
}

func (bh *BookingHandler) CheckoutWithSnap() echo.HandlerFunc {
	return func(c echo.Context) error {

		var input = new(CheckoutRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		res, mt, err := bh.s.CheckoutWithSnap(input.BookingID)

		if err != nil {
			c.Logger().Error("Handler: Input Process Error (Checkout): ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		fmt.Println("Res payment", res)
		fmt.Println("Res midtrans", mt)

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "Success checkout", mt))
	}
}

func (bh *BookingHandler) GetAllBookingDemo() echo.HandlerFunc {
	return func(c echo.Context) error {
		statusBooking, err := strconv.Atoi(c.QueryParam("booking_status"))

		if err != nil {
			statusBooking = 0
		}

		search, page, pageSize, err := helper.GetPaginationQuery(c)

		if err != nil {
			c.Logger().Error("Handler : Get Query Param Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Query Param Error", nil))
		}

		result, totalPage, totalItems, err := bh.s.GetAllBookingDemo(*search, *page, *pageSize, uint(statusBooking))

		if len(result) == 0 {
			return c.JSON(http.StatusOK, helper.FormatResponse(true, "Data is empty", nil))
		}

		if err != nil {
			c.Logger().Error("Handler : Get All Booking Demo Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get All Booking Demo Error", nil))
		}

		var allResult = new(GetResponse)
		allResult.MainData = result
		allResult.Page = page
		allResult.PageSize = pageSize
		allResult.TotalPage = totalPage
		allResult.TotalItems = totalItems

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get All Data Booking Demo", allResult))
	}
}

func (bh *BookingHandler) GetAllBooking() echo.HandlerFunc {
	return func(c echo.Context) error {
		statusBooking, err := strconv.Atoi(c.QueryParam("booking_status"))

		if err != nil {
			statusBooking = 0
		}

		search, page, pageSize, err := helper.GetPaginationQuery(c)

		if err != nil {
			c.Logger().Error("Handler : Get Query Param Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Query Param Error", nil))
		}

		result, totalPage, totalItems, err := bh.s.GetAllBooking(*search, *page, *pageSize, uint(statusBooking))

		if len(result) == 0 {
			return c.JSON(http.StatusOK, helper.FormatResponse(true, "Data is empty", nil))
		}

		if err != nil {
			c.Logger().Error("Handler : Get All Booking Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get All Booking Error", nil))
		}

		var allResult = new(GetResponse)
		allResult.MainData = result
		allResult.Page = page
		allResult.PageSize = pageSize
		allResult.TotalPage = totalPage
		allResult.TotalItems = totalItems

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get All Data Booking Demo", allResult))
	}
}

func (bh *BookingHandler) GetCustomersBooking() echo.HandlerFunc {
	return func(c echo.Context) error {

		qryParam := c.QueryParams()

		statusBooking, err := strconv.Atoi(qryParam.Get("booking_status"))

		if err != nil {
			statusBooking = 0
		}

		userID, err := bh.jwt.GetID(c)

		if err != nil {
			var result = []any{}
			return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success", result))
		}

		result, err := bh.s.GetCustomersBooking(userID, uint(statusBooking))

		if err != nil {
			c.Logger().Error("Handler : Get Booking By ID Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Booking By ID Error", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get Booking By ID Data", result))
	}

}

func (bh *BookingHandler) GetCustomersBookingByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramID = c.Param("id")
		// id, err := strconv.Atoi(paramID)
		// if err != nil {
		// 	c.Logger().Error("Handler : Param ID Error : ", err.Error())
		// 	return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid User Input Param ID", nil))
		// }

		err := uuid.Validate(paramID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Please enter valid booking id", nil))
		}

		id := uuid.MustParse(paramID)

		result, err := bh.s.GetCustomersBookingByID(id)

		if err != nil {
			c.Logger().Error("Handler : Get Booking By ID Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Booking By ID Error", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get Booking By ID Data", result))
	}

}

func (bh *BookingHandler) GetBookingByBookingInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramID = c.Param("booking_invoice")

		result, err := bh.s.GetBookingByBookingInvoice(paramID)

		if err != nil {
			c.Logger().Error("Handler : Get Booking By ID Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Booking By ID Error", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get Booking By ID Data", result))
	}

}

func (bh *BookingHandler) GetBookingByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramID = c.Param("id")
		id := uuid.MustParse(paramID)

		result, err := bh.s.GetBookingByID(id)

		if err != nil {
			c.Logger().Error("Handler : Get Booking By ID Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Get Booking By ID Error", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Get Booking By ID Data", result))
	}

}

func (bh *BookingHandler) UpdateBooking() echo.HandlerFunc {
	return func(c echo.Context) error {
		var paramID = c.Param("id")

		id := uuid.MustParse(paramID)

		var input = new(InputRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		// res, err := bh.s.GetPaymentByID(input.PaymentType)

		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "payment not found", nil))
		// }

		var serviceUpdateBooking = new(booking.Booking)

		serviceUpdateBooking.OfficeID = input.OfficeID
		// serviceUpdateBooking.PaymentType = res.Name
		serviceUpdateBooking.BookingStartDate = input.BookingStartDate
		serviceUpdateBooking.BookingEndDate = input.BookingEndDate
		serviceUpdateBooking.UsesDescription = input.UsesDescription

		resultUpdatedBooking, name, err := bh.s.UpdateBooking(id, *serviceUpdateBooking)

		if err != nil {
			c.Logger().Error("Handler: Input Process Error (UpdateBooking): ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		var response = new(UpdateResponse)
		response.CustomerName = *name
		response.HourCount = resultUpdatedBooking.HourCount
		response.DaysCount = resultUpdatedBooking.DaysCount
		response.PaymentType = resultUpdatedBooking.PaymentType
		response.Discount = resultUpdatedBooking.Discount
		response.BookingStatus = resultUpdatedBooking.BookingStatus
		response.BookingStartDate = resultUpdatedBooking.BookingStartDate
		response.BookingEndDate = resultUpdatedBooking.BookingEndDate
		response.UsesDescription = resultUpdatedBooking.UsesDescription
		response.TotalPrice = resultUpdatedBooking.TotalPrice

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success update Booking", response))
	}
}

func (bh *BookingHandler) DeleteBooking() echo.HandlerFunc {
	return func(c echo.Context) error {
		paramID := c.Param("id")
		// id, err := strconv.Atoi(paramID)
		// if err != nil {
		// 	c.Logger().Error("Handler : Param ID Error : ", err.Error())
		// 	return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Invalid User Input Param ID", nil))
		// }
		id := uuid.MustParse(paramID)

		result, err := bh.s.DeleteBooking(id)
		if err != nil {
			c.Logger().Error("Handler : Delete Booking By ID Error : ", err.Error())
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Delete Booking By ID Error: "+err.Error(), nil))
		}

		if !result {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Delete Booking By ID Failed", nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success Delete Booking By ID Data", nil))
	}
}

func (bh *BookingHandler) ReportBooking() echo.HandlerFunc {
	return func(c echo.Context) error {
		bookings, err := bh.s.ReportBooking()

		if err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		var reportData []BookingReport

		for _, booking := range bookings {
			reportData = append(reportData, BookingReport{
				ID:               booking.ID,
				OfficeName:       booking.OfficeName,
				CustomerName:     booking.CustomerName,
				OfficeID:         booking.OfficeID,
				UserID:           booking.UserID,
				CustomerID:       booking.CustomerID,
				TotalPrice:       booking.TotalPrice,
				HourCount:        booking.HourCount,
				DaysCount:        booking.DaysCount,
				WeekCount:        booking.WeekCount,
				MonthCount:       booking.MonthCount,
				PaymentType:      booking.PaymentType,
				Discount:         booking.Discount,
				BookingStatus:    booking.BookingStatus,
				BookingStartDate: booking.BookingStartDate,
				BookingEndDate:   booking.BookingEndDate,
				UsesDescription:  booking.UsesDescription,
				VAAccount:        booking.VAAccount,
				Notes:            booking.Notes,
				Payment:          booking.Payment,
			})
		}

		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		index, err := f.NewSheet("Bookings Report")
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Failed to create report", nil))
		}
		header := []string{"Booking ID", "Office Name", "Customer Name", "Office ID", "User ID", "Customer ID", "Total Price", "Hour Count", "Days Count", "Halfday Count", "Week Count", "Month Count", "Payment Type", "Discount", "Booking Status", "Booking Start Date", "Booking End Date", "Uses Description", "VA Account", "Notes", "Payment"}
		f.SetSheetRow("Bookings Report", "A1", &header)

		row := 2
		for _, report := range reportData {
			f.SetCellValue("Bookings Report", fmt.Sprintf("A%d", row), report.ID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("B%d", row), report.OfficeName)
			f.SetCellValue("Bookings Report", fmt.Sprintf("C%d", row), report.CustomerName)
			f.SetCellValue("Bookings Report", fmt.Sprintf("D%d", row), report.OfficeID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("E%d", row), report.UserID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("F%d", row), report.CustomerID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("G%d", row), report.TotalPrice)
			f.SetCellValue("Bookings Report", fmt.Sprintf("H%d", row), report.HourCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("I%d", row), report.DaysCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("J%d", row), report.HalfDayCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("K%d", row), report.WeekCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("L%d", row), report.MonthCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("M%d", row), report.PaymentType)
			f.SetCellValue("Bookings Report", fmt.Sprintf("N%d", row), report.Discount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("O%d", row), report.BookingStatus)
			f.SetCellValue("Bookings Report", fmt.Sprintf("P%d", row), report.BookingStartDate)
			f.SetCellValue("Bookings Report", fmt.Sprintf("Q%d", row), report.BookingEndDate)
			f.SetCellValue("Bookings Report", fmt.Sprintf("R%d", row), report.UsesDescription)
			f.SetCellValue("Bookings Report", fmt.Sprintf("S%d", row), report.VAAccount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("T%d", row), report.Notes)
			f.SetCellValue("Bookings Report", fmt.Sprintf("U%d", row), report.Payment)
			row++
		}

		now := time.Now()

		fileName := fmt.Sprintf("REPORT_BOOKING_%s.xlsx", now.Format("2006-01-02-15-04-05"))
		currentDir, err := os.Getwd()
		// // fmt.Println(currentDir) debugging
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error getting working directory", nil))
		}

		basePath := filepath.Join(currentDir, "/assets/data/document")
		filePath := filepath.Join(basePath, fileName)

		f.SetActiveSheet(index)
		if err := f.SaveAs(filePath); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Error creating the file", nil))
		}

		return c.Inline(filePath, fileName)
	}
}

func (bh *BookingHandler) ReportBookingCustomer() echo.HandlerFunc {
	return func(c echo.Context) error {
		bookings, err := bh.s.ReportBooking()

		if err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		var reportData []BookingReport

		for _, booking := range bookings {
			reportData = append(reportData, BookingReport{
				ID:               booking.ID,
				OfficeName:       booking.OfficeName,
				CustomerName:     booking.CustomerName,
				OfficeID:         booking.OfficeID,
				UserID:           booking.UserID,
				CustomerID:       booking.CustomerID,
				TotalPrice:       booking.TotalPrice,
				HourCount:        booking.HourCount,
				DaysCount:        booking.DaysCount,
				WeekCount:        booking.WeekCount,
				MonthCount:       booking.MonthCount,
				PaymentType:      booking.PaymentType,
				Discount:         booking.Discount,
				BookingStatus:    booking.BookingStatus,
				BookingStartDate: booking.BookingStartDate,
				BookingEndDate:   booking.BookingEndDate,
				UsesDescription:  booking.UsesDescription,
				VAAccount:        booking.VAAccount,
				Notes:            booking.Notes,
				Payment:          booking.Payment,
			})
		}

		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()

		index, err := f.NewSheet("Bookings Report")
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Failed to create report", nil))
		}
		header := []string{"Booking ID", "Office Name", "Customer Name", "Office ID", "User ID", "Customer ID", "Total Price", "Hour Count", "Days Count", "Halfday Count", "Week Count", "Month Count", "Payment Type", "Discount", "Booking Status", "Booking Start Date", "Booking End Date", "Uses Description", "VA Account", "Notes", "Payment"}
		f.SetSheetRow("Bookings Report", "A1", &header)

		row := 2
		for _, report := range reportData {
			f.SetCellValue("Bookings Report", fmt.Sprintf("A%d", row), report.ID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("B%d", row), report.OfficeName)
			f.SetCellValue("Bookings Report", fmt.Sprintf("C%d", row), report.CustomerName)
			f.SetCellValue("Bookings Report", fmt.Sprintf("D%d", row), report.OfficeID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("E%d", row), report.UserID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("F%d", row), report.CustomerID)
			f.SetCellValue("Bookings Report", fmt.Sprintf("G%d", row), report.TotalPrice)
			f.SetCellValue("Bookings Report", fmt.Sprintf("H%d", row), report.HourCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("I%d", row), report.DaysCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("J%d", row), report.HalfDayCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("K%d", row), report.WeekCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("L%d", row), report.MonthCount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("M%d", row), report.PaymentType)
			f.SetCellValue("Bookings Report", fmt.Sprintf("N%d", row), report.Discount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("O%d", row), report.BookingStatus)
			f.SetCellValue("Bookings Report", fmt.Sprintf("P%d", row), report.BookingStartDate)
			f.SetCellValue("Bookings Report", fmt.Sprintf("Q%d", row), report.BookingEndDate)
			f.SetCellValue("Bookings Report", fmt.Sprintf("R%d", row), report.UsesDescription)
			f.SetCellValue("Bookings Report", fmt.Sprintf("S%d", row), report.VAAccount)
			f.SetCellValue("Bookings Report", fmt.Sprintf("T%d", row), report.Notes)
			f.SetCellValue("Bookings Report", fmt.Sprintf("U%d", row), report.Payment)
			row++
		}

		now := time.Now()

		fileName := fmt.Sprintf("REPORT_%s.xlsx", now.Format("2006-01-02-15-04-05"))
		currentDir, err := os.Getwd()
		// // fmt.Println(currentDir) debugging
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, "Error getting working directory", nil))
		}

		basePath := filepath.Join(currentDir, "/assets/data/document")
		filePath := filepath.Join(basePath, fileName)

		f.SetActiveSheet(index)
		if err := f.SaveAs(filePath); err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Error creating the file", nil))
		}

		return c.Inline(filePath, fileName)
	}

}

func (bh *BookingHandler) CancelBooking() echo.HandlerFunc {
	return func(c echo.Context) error {

		getID, _ := bh.jwt.GetID(c)

		var paramID = c.Param("id")

		id := uuid.MustParse(paramID)

		var input = new(CancelRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		resBook, err := bh.s.GetCustomersBookingByID(id)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		if resBook.UserID != getID {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "this booking do not belong to the authenticated user", nil))

		}

		_, err = bh.s.CancelBooking(id, input.Explanation)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success update booking", nil))
	}
}

func (bh *BookingHandler) CancelBookingByAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {

		role := bh.jwt.CheckRole(c)

		if role != "SUPERADMIN" {
			return c.JSON(http.StatusUnauthorized, helper.FormatResponse(false, "you are not to allow to use this feature", nil))
		}

		var paramID = c.Param("id")

		id := uuid.MustParse(paramID)

		var input = new(CancelRequest)

		if err := c.Bind(input); err != nil {
			c.Logger().Error("Handler : Bind Input Error : ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Bind input Error", nil))
		}

		_, err := bh.s.CancelBooking(id, input.Explanation)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "Success update booking", nil))
	}
}

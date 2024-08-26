package main

import (
	configs "ceo-suite-go/configs"
	"ceo-suite-go/helper"
	"ceo-suite-go/helper/email"
	encrypt "ceo-suite-go/helper/encrypt"
	"ceo-suite-go/routes"
	"ceo-suite-go/utils/watcher"
	"fmt"
	"net/http"

	"ceo-suite-go/utils/database"
	"ceo-suite-go/utils/midtrans"

	dataBooking "ceo-suite-go/features/booking/data"
	handlerBooking "ceo-suite-go/features/booking/handler"
	serviceBooking "ceo-suite-go/features/booking/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {

	e := echo.New()

	e.Static("/static", "assets")
	var config = configs.InitConfig()

	var midtrans = midtrans.InitMidtrans(*config)

	db, err := database.InitDB(*config)
	if err != nil {
		e.Logger.Fatal("Cannot run database: ", err.Error())
	}
	var storage = helper.InitStorage(*config)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	e.GET("/api/v1", func(c echo.Context) error {
		return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "Endpoint not found", nil))
	})

	var encrypt = encrypt.New()
	var email = email.New(*config)

	jwtInterface := helper.New(config.Secret, config.RefSecret)

	bookingModel := dataBooking.New(db, customerModel, officeModel)

	bookingServices := serviceBooking.New(bookingModel, userServices, jwtInterface, email, midtrans, customerService)

	bookingController := handlerBooking.NewHandler(bookingServices, jwtInterface, officeServices)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	group := e.Group("/api/v1")

	routes.RouteBooking(group, bookingController, *config)


	wtc := watcher.New(db, bookingModel, userModel, *config)
	err = wtc.DoEveryTenMinute()

	if err != nil {
		logrus.Info("[WATCHER] Watcher is not active")
	}

	e.Logger.Debug(db)

	e.Logger.Info(fmt.Sprintf("Listening in port :%d", config.ServerPort))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}

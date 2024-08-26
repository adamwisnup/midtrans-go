package seeder

import (
	custData "ceo-suite-go/features/customer/data"
	officeData "ceo-suite-go/features/office/data"
	userData "ceo-suite-go/features/users/data"
	models "ceo-suite-go/utils/database/seeder/models"
	seeder "ceo-suite-go/utils/database/seeder/type"
	"time"

	"gorm.io/gorm"
)

func All() []seeder.Seed {

	isTrue := true
	// config := configs.ReadData()
	var seeds []seeder.Seed = []seeder.Seed{
		// {
		// 	Name: "CreateRegion1",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "ACEH")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion2",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SUMATERA UTARA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion3",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SUMATERA BARAT")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion4",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "RIAU")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion5",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "JAMBI")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion6",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SUMATERA SELATAN")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion7",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "BENGKULU")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion8",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "LAMPUNG")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion9",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KEPULAUAN BANGKA BELITUNG")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion10",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KEPULAUAN RIAU")
		// 	},
		// },
		{
			Name: "CreateRegion",
			Run: func(db *gorm.DB) error {
				return models.CreateRegions(db, "Jakarta, Sahid Sudirman Center")
			},
		},
		{
			Name: "CreateRegion2",
			Run: func(db *gorm.DB) error {
				return models.CreateRegions(db, "Jakarta, One Pacific Place")
			},
		},
		{
			Name: "CreateRegion3",
			Run: func(db *gorm.DB) error {
				return models.CreateRegions(db, "Jakarta, Wisma GKBI")
			},
		},
		{
			Name: "CreateRegion4",
			Run: func(db *gorm.DB) error {
				return models.CreateRegions(db, "Jakarta, AXA Tower")
			},
		},
		{
			Name: "CreateRegion5",
			Run: func(db *gorm.DB) error {
				return models.CreateRegions(db, "Jakarta, Indonesia Stock Exchange")
			},
		},
		// {
		// 	Name: "CreateRegion12",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "JAWA BARAT")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion13",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "JAWA TENGAH")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion14",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "DAERAH ISTIMEWA YOGYAKARTA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion15",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "JAWA TIMUR")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion16",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "BANTEN")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion17",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "BALI")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion18",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "NUSA TENGGARA BARAT")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion19",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "NUSA TENGGARA TIMUR")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion20",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KALIMANTAN BARAT")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion21",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KALIMANTAN TENGAH")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion22",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KALIMANTAN SELATAN")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion23",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KALIMANTAN TIMUR")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion24",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "KALIMANTAN UTARA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion25",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SULAWESI UTARA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion26",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SULAWESI TENGAH")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion27",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SULAWESI SELATAN")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion28",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SULAWESI TENGGARA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion29",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "GORONTALO")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion30",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "SULAWESI BARAT")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion31",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "MALUKU")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion32",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "MALUKU UTARA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion33",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "PAPUA")
		// 	},
		// },
		// {
		// 	Name: "CreateRegion34",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateRegions(db, "PAPUA BARAT")
		// 	},
		// },

		{
			Name: "Create Admin",
			Run: func(db *gorm.DB) error {
				return models.CreateUsers(db, userData.User{
					Name:        "Admin",
					Email:       "admin@admin.com",
					Password:    "admin123!",
					Role:        "SUPERADMIN",
					DateOfBirth: time.Now(),
					PhoneNumber: "087755782234",
					// TokenResetPass: "",
					Status: "ACTIVE",
					// Avatar: "https://cdn2.iconfinder.com/data/icons/user-interface-essential-solid/32/Artboard_52-1024.png",
				})
			},
		},

		{
			Name: "Create Customer",
			Run: func(db *gorm.DB) error {
				return models.CreateUsers(db, userData.User{
					Name:        "Customer",
					Email:       "customer@customer.com",
					Password:    "customer123!",
					Role:        "CUSTOMER",
					DateOfBirth: time.Now(),
					PhoneNumber: "087755782234",
					// TokenResetPass: "",
					Status: "ACTIVE",
					// Avatar: "https://cdn2.iconfinder.com/data/icons/user-interface-essential-solid/32/Artboard_52-1024.png",
				})
			},
		},

		{
			Name: "Create Customer Details",
			Run: func(db *gorm.DB) error {
				return models.CreateCustomer(db, custData.Customer{
					UserID:         2,
					FullName:       "Customer customer",
					Position:       "CEO",
					CompanyName:    "Customer Sdn. Bhd.",
					CompanyEmail:   "customer@customer.com",
					CompanyAddress: "Customer Address",
				})
			},
		},

		{
			Name: "Create payment data",
			Run: func(d *gorm.DB) error {
				return models.CreatePaymentList(d, "BCA")
			},
		},

		{
			Name: "Create payment data 2",
			Run: func(d *gorm.DB) error {
				return models.CreatePaymentList(d, "BNI")
			},
		},

		// {
		// 	Name: "Create payment data 3",
		// 	Run: func(d *gorm.DB) error {
		// 		return models.CreatePaymentList(d, "BRI")
		// 	},
		// },

		// {
		// 	Name: "Create payment data 4",
		// 	Run: func(d *gorm.DB) error {
		// 		return models.CreatePaymentList(d, "QRIS")
		// 	},
		// },

		{
			Name: "Create payment data 5",
			Run: func(d *gorm.DB) error {
				return models.CreatePaymentList(d, "GOPAY")
			},
		},
		// {
		// 	Name: "CreateOffice1",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateOffice(db, data.Office{
		// 			Name:           "Axa Boardroom",
		// 			Discount:       0,
		// 			OfficeCategory: 1,
		// 			Description:    "Axa Boardroom Description",
		// 			Status:         true,
		// 			RegionID:       11,
		// 			UserID:         1,
		// 		}, data.Price{
		// 			PricePerHour: 480000,
		// 			PriceDaily:   2880000,
		// 			PriceHalfDay: 1920000,
		// 		}, []data.OfficeCatalogue{
		// 			{
		// 				URL: "https://ceosuite.com/wp-content/uploads/2020/07/AXA-Boardroom-1-1024x683.jpg",
		// 			},
		// 		})
		// 	},
		// },
		{
			Name: "CreateOfficeCategory 1",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Office")
			},
		},
		{
			Name: "CreateOfficeCategory 2",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Meeting Room")
			},
		},

		{
			Name: "CreateOfficeCategory 3",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Board Room")
			},
		},

		{
			Name: "CreateOfficeCategory 4",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Class Room")
			},
		},

		{
			Name: "CreateOfficeCategory 5",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Co Working Space Dedicated")
			},
		},

		{
			Name: "CreateOfficeCategory 6",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Co Working Space Non Dedicated")
			},
		},

		{
			Name: "CreateOfficeCategory 7",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Seminar Room")
			},
		},

		{
			Name: "CreateOfficeCategory 8",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Theathre Room")
			},
		},

		{
			Name: "CreateOfficeCategory 9",
			Run: func(db *gorm.DB) error {
				return models.CreateOfficeCategory(db, "Combined Boardroom Meeting Room")
			},
		},

		{
			Name: "Create Office Example Data",
			Run: func(db *gorm.DB) error {
				return models.CreateOffice(db, officeData.Office{
					Name:           "Example Office",
					OfficeCategory: 3,
					OfficeNumber:   "OFF.XB.NT.03",
					Description:    "Example Office Description",
					Lat:            -6.2239292,
					Lng:            106.806077,
					Status:         &isTrue,
					RegionID:       11,
					UserID:         1,
				}, officeData.Price{
					PricePerHour: 1000000,
					PriceHalfDay: 10000000,
					PriceDaily:   15000000,
					PriceWeekly:  40000000,
					PriceMonthly: 100000000,
					PriceMinimum: 2000000,
				}, officeData.OfficeDetails{
					Location:       "Example Location Office Street",
					ServiceName:    "Example Location Office Service Name",
					CapacityPerson: 50,
					CapacityDesk:   25,
					CapacityChair:  25,
					Address:        "Example Location Office Address",
				}, []officeData.OfficeCatalogue{
					{
						URL: "https://ceosuite-api.nikici.com/static/office/images/idx/pantry/IDX_Pantry.jpeg",
					},
				})
			},
		},

		// {
		// 	Name: "CreateOffice2",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateOffice(db, data.Office{
		// 			Name:        "GKBI Office",
		// 			RoomType:    "Room A",
		// 			Description: "Stay inspired in a shared environment with like-minded individuals pursuing their passions, and become part of a dynamic creative community that stretches throughout Asia.",
		// 			ImagePath:   fmt.Sprintf("%s/static/office/images/gkbi-office.png", config.BaseURL),
		// 			Capacity:    10,
		// 			Table:       5,
		// 			Chair:       5,
		// 			RegionID:    11,
		// 		})
		// 	},
		// },
		// {
		// 	Name: "CreateOffice3",
		// 	Run: func(db *gorm.DB) error {
		// 		return models.CreateOffice(db, data.Office{
		// 			Name:        "IDX Office",
		// 			RoomType:    "Room A",
		// 			Description: "Stay inspired in a shared environment with like-minded individuals pursuing their passions, and become part of a dynamic creative community that stretches throughout Asia.",
		// 			ImagePath:   fmt.Sprintf("%s/static/office/images/IDX-Suite-07.jpg", config.BaseURL),
		// 			Capacity:    10,
		// 			Table:       5,
		// 			Chair:       5,
		// 			RegionID:    11,
		// 		})
		// 	},
		// },
	}
	return seeds
}

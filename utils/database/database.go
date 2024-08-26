package database

import (
	configs "ceo-suite-go/configs"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(c configs.ProgrammingConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost,
		c.DBUser,
		c.DBPass,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("Asia/Jakarta")
			return time.Now().In(utc)
		},
	})

	if err != nil {
		log.Error("Terjadi kesalahan pada database, error:", err.Error())
		return nil, err
	}

	return db, nil
}

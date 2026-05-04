package db

import (
	"fmt"
	"os"

	"github.com/xerox-1315/TravelShare.git/errs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() (*gorm.DB, error) {
	// строка подключения к БД, данные берем из файле с перемнными окружения
	connect_str := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	// подключаемся к БД
	db, err := gorm.Open(postgres.Open(connect_str), &gorm.Config{})
	if err != nil {
		return nil, errs.ErrorDBConnection
	}
	return db, nil
}

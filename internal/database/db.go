package database

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dns := fmt.Sprintf("user=%s password=%s  host=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADDRESS"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	_db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("initial database error: %s", err))
	}

	_db.AutoMigrate(&UserInfo{})

	db = _db
}

func GetDB() *gorm.DB {
	if nil == db {
		panic(types.Error{Msg: "should init database before get it."})
	}
	return db
}

func DestroyDB() {
	if db == nil {
		return
	}
}

package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Sửa thông tin kết nối theo DB của bạn: user:password@tcp(127.0.0.1:3306)/wedding
	dsn := "wedding_user:cuoihoatuyet@tcp(103.216.119.197:3307)/wedding?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Giới hạn connection pool: mặc định database/sql không giới hạn số kết nối mở,
	// dễ tốn RAM/CPU (cả phía MySQL) khi có nhiều request cùng lúc.
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	fmt.Println("Database connection successfully established!")
	DB = database
}

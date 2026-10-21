package config

import (
	"fmt"
	"log"

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

	fmt.Println("Database connection successfully established!")
	DB = database
}

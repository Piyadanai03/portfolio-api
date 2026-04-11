package config

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ ไม่สามารถเชื่อมต่อ Database ได้: ", err)
	}

	log.Println("✅ เชื่อมต่อ Database สำเร็จ!")
	DB = database
};
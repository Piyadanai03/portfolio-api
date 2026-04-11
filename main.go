package main

import (
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/routes"
	"github.com/joho/godotenv"
	"os"
)

// @title           Piyadanai Portfolio API
// @version         1.0
// @description     นี่คือระบบ Backend สำหรับ Portfolio ของคุณปิยดนัย
// @termsOfService  http://swagger.io/terms/

// @contact.name   Piyadanai Krongklang
// @contact.url    https://github.com/Piyadanai03

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description ใส่คำว่า "Bearer " ตามด้วย Token ของคุณ
func main() {
	godotenv.Load()
	config.ConnectDatabase()
	// models.MigrateDB(config.DB)

	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

package profile

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// ดึงข้อมูล Profile (พร้อม Contacts)
func GetProfile(c *gin.Context) {
	// ดึง user_id มาจาก JWT Token (Middleware ใส่ไว้ให้แล้ว)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	// ใช้ Preload("Contacts") เพื่อให้มันดึงตาราง Contact ติดมาด้วย
	if err := config.DB.Preload("Contacts").Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้"})
		return
	}

	c.JSON(http.StatusOK, user)
}
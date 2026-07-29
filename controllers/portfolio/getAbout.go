package portfolio

import (
	"net/http"
	"os"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

func GetAbout(c *gin.Context) {
	var user models.User
	var achievements []models.Achievement

	ownerID := os.Getenv("USER_ID")

	if ownerID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ระบบยังไม่ได้ตั้งค่า USER_ID ของเจ้าของเว็บ"})
		return
	}

	// Fetch user with all related data
	if err := config.DB.
		Preload("Experiences").
		Preload("Studies").
		Preload("Contacts").
		Preload("Projects").
		Where("id = ?", ownerID).
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลโปรไฟล์ของเจ้าของเว็บ"})
		return
	}

	// Fetch achievements for the user
	if err := config.DB.
		Where("user_id = ?", ownerID).
		Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูล achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"achievements": achievements,
	})
}

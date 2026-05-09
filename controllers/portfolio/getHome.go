package portfolio

import (
	"net/http"
	"os"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

func GetHomeData(c *gin.Context) {
	var user models.User

	ownerID := os.Getenv("USER_ID")
	
	if ownerID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ระบบยังไม่ได้ตั้งค่า USER_ID ของเจ้าของเว็บ"})
		return
	}

	if err := config.DB.Preload("Projects").Where("id = ?", ownerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลโปรไฟล์ของเจ้าของเว็บ"})
		return
	}

	c.JSON(http.StatusOK, user)
}
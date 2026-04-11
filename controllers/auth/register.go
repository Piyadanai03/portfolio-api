package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		FullName string `json:"full_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- ส่วนที่เพิ่มเข้ามา: การเข้ารหัสรหัสผ่าน ---
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเข้ารหัสรหัสผ่านได้"})
		return
	}
	// ---------------------------------------

	user := models.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword), // เซฟตัวที่ถูก Hash แล้วลงไป
		FullName:     input.FullName,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username นี้ถูกใช้ไปแล้ว"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลงทะเบียน Admin แบบปลอดภัยสำเร็จ!"})
}
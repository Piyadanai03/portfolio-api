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

// อัปเดตข้อมูล Profile
func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 1. รับข้อมูล JSON จาก React
	var input struct {
		FullName string           `json:"fullName"`
		Position string           `json:"position"`
		BioText  string           `json:"bioText"`
		Address  string           `json:"address"`
		ProfileImageURL string    `json:"profileImageURL"`
		ResumeURL       string    `json:"resumeURL"`
		Contacts []models.Contact `json:"contacts"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูลไม่ถูกต้อง", "details": err.Error()})
		return
	}

	// 2. ค้นหา User คนปัจจุบัน
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้"})
		return
	}

	// 3. อัปเดตข้อมูล Text ทั่วไป
	user.FullName = input.FullName
	user.Position = input.Position
	user.BioText = input.BioText
	user.Address = input.Address
	user.ProfileImageURL = input.ProfileImageURL
	user.ResumeURL = input.ResumeURL
	config.DB.Save(&user)

	// 4. พระเอก: อัปเดต Contacts ทั้งหมดด้วยท่า Association Replace
	// ท่านี้ GORM จะฉลาดมาก: มันจะเช็คว่าอันไหนถูกลบ มันจะลบให้ อันไหนของใหม่ มันจะสร้างให้!
	err := config.DB.Model(&user).Association("Contacts").Replace(input.Contacts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตช่องทางติดต่อได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตโปรไฟล์สำเร็จ!"})
}
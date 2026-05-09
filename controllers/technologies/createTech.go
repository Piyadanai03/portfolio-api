package technologies

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

func CreateTech(c *gin.Context) {
	// รับข้อมูลจาก Form-Data แบบพื้นฐาน
	name := c.PostForm("name")
	category := c.PostForm("category")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุชื่อเทคโนโลยี (name)"})
		return
	}

	// จัดการอัปโหลด Icon Image
	var iconURL string
	file, _, err := c.Request.FormFile("icon")
	if err == nil {
		defer file.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_technologies")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดรูปไอคอนล้มเหลว"})
			return
		}
		iconURL = uploadedURL
	}

	// สร้าง Object เทคโนโลยี
	tech := models.Technology{
		Name:     name,
		Category: category,
		IconURL:  iconURL,
	}

	// บันทึกเทคโนโลยีลงฐานข้อมูล
	if err := config.DB.Create(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างเทคโนโลยีไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เทคโนโลยีถูกสร้างเรียบร้อยแล้ว", "technology": tech})
}
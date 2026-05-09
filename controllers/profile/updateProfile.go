package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils" // 🌟 Import โฟลเดอร์ utils
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateProfile อัปเดตข้อมูล Profile แบบครบวงจร
func UpdateProfile(c *gin.Context) {
	// 1. ดึง userID จาก Token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลยืนยันตัวตน"})
		return
	}

	// 2. ค้นหา User คนปัจจุบัน
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้"})
		return
	}

	// 3. อ่านข้อมูล Text
	fullName := c.PostForm("fullName")
	position := c.PostForm("position")
	bioText := c.PostForm("bioText")
	address := c.PostForm("address")
	contactsJSON := c.PostForm("contacts")

	// 4. อัปโหลดรูปภาพ Profile (ถ้ามี)
	profileImageURL := user.ProfileImageURL
	profileFile, _, err := c.Request.FormFile("profileImage")
	if err == nil {
		defer profileFile.Close()
		// 🌟 เรียกใช้จาก utils
		uploadedURL, uploadErr := utils.UploadToCloudinary(profileFile, "portfolio_profiles")
		if uploadErr == nil {
			profileImageURL = uploadedURL
		}
	}

	// 5. อัปโหลด Resume (ถ้ามี)
	resumeURL := user.ResumeURL
	resumeFile, _, err := c.Request.FormFile("resume")
	if err == nil {
		defer resumeFile.Close()
		// 🌟 เรียกใช้จาก utils
		uploadedURL, uploadErr := utils.UploadToCloudinary(resumeFile, "portfolio_resumes")
		if uploadErr == nil {
			resumeURL = uploadedURL
		}
	}

	// 6. บันทึกข้อมูลส่วนตัวลง DB
	if fullName != "" {
		user.FullName = fullName
	}
	if position != "" {
		user.Position = position
	}
	if bioText != "" {
		user.BioText = bioText
	}
	if address != "" {
		user.Address = address
	}
	user.ProfileImageURL = profileImageURL
	user.ResumeURL = resumeURL

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถบันทึกข้อมูลผู้ใช้ได้"})
		return
	}

	// 7. จัดการข้อมูล Contacts ด้วย Transaction
	if contactsJSON != "" {
		var contacts []models.Contact
		if err := json.Unmarshal([]byte(contactsJSON), &contacts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูล Contacts ไม่ถูกต้อง"})
			return
		}

		err := config.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("user_id = ?", user.ID).Delete(&models.Contact{}).Error; err != nil {
				return err
			}
			if len(contacts) > 0 {
				for i := range contacts {
					contacts[i].UserID = user.ID
				}
				if err := tx.Create(&contacts).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			fmt.Println("❌ Contacts Update Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตช่องทางติดต่อได้"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "อัปเดตโปรไฟล์สำเร็จ!",
		"profileImageURL": profileImageURL,
		"resumeURL":       resumeURL,
	})
}
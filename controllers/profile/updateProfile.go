package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateProfile อัปเดตข้อมูล Profile แบบครบวงจร
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลยืนยันตัวตน"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้"})
		return
	}

	fullName := c.PostForm("fullName")
	position := c.PostForm("position")
	bioText := c.PostForm("bioText")
	address := c.PostForm("address")
	contactsJSON := c.PostForm("contacts")

	// 🌟 1. รับค่า URL ที่พิมพ์มาจากฟอร์ม ถ้ามีค่าให้เอามาใช้ก่อน
	profileImageURL := c.PostForm("profileImageURL")
	if profileImageURL == "" {
		profileImageURL = user.ProfileImageURL
	}

	resumeURL := c.PostForm("resumeURL")
	if resumeURL == "" {
		resumeURL = user.ResumeURL
	}

	// 🌟 2. อัปโหลดรูปภาพ Profile (ถ้ามีไฟล์แนบมา จะทับ URL ด้านบนเสมอ)
	profileFile, _, err := c.Request.FormFile("profileImage")
	if err == nil {
		defer profileFile.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(profileFile, "portfolio_profiles")
		if uploadErr == nil {
			profileImageURL = uploadedURL
		}
	}

	// 🌟 3. อัปโหลด Resume (ถ้ามีไฟล์แนบมา จะทับ URL ด้านบนเสมอ)
	resumeFile, _, err := c.Request.FormFile("resume")
	if err == nil {
		defer resumeFile.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(resumeFile, "portfolio_resumes")
		if uploadErr == nil {
			resumeURL = uploadedURL
		}
	}

	// 4. บันทึกข้อมูลส่วนตัวลง DB
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

	// 🌟 5. จัดการข้อมูล Contacts และอัปโหลด Icon ขึ้น Cloudinary
	if contactsJSON != "" {
		var contacts []models.Contact
		if err := json.Unmarshal([]byte(contactsJSON), &contacts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูล Contacts ไม่ถูกต้อง"})
			return
		}

		// ตรวจสอบและอัปโหลดรูป Icon ที่เป็น Base64 ขึ้น Cloudinary
		for i := range contacts {
			if strings.HasPrefix(contacts[i].IconURL, "data:image") {
				uploadedURL, uploadErr := utils.UploadToCloudinary(contacts[i].IconURL, "portfolio_contact_icons")
				if uploadErr == nil {
					contacts[i].IconURL = uploadedURL
				} else {
					fmt.Println("⚠️ Failed to upload contact icon:", uploadErr)
				}
			}
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
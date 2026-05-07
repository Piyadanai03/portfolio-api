package profile

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

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

	profileImageURL := user.ProfileImageURL
	profileFile, _, err := c.Request.FormFile("profileImage")
	if err == nil {
		defer profileFile.Close()
		uploadedURL, uploadErr := uploadToCloudinary(profileFile, "portfolio_profiles")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดรูปโปรไฟล์ล้มเหลว: " + uploadErr.Error()})
			return
		}
		profileImageURL = uploadedURL
	}

	resumeURL := user.ResumeURL
	resumeFile, _, err := c.Request.FormFile("resume")
	if err == nil {
		defer resumeFile.Close()
		uploadedURL, uploadErr := uploadToCloudinary(resumeFile, "portfolio_resumes")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลด Resume ล้มเหลว: " + uploadErr.Error()})
			return
		}
		resumeURL = uploadedURL
	}

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
	config.DB.Save(&user)

	if contactsJSON != "" {
		var contacts []models.Contact
		if err := json.Unmarshal([]byte(contactsJSON), &contacts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบข้อมูล Contacts ไม่ถูกต้อง"})
			return
		}
		if err := config.DB.Model(&user).Association("Contacts").Replace(contacts); err != nil {
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

func uploadToCloudinary(file interface{}, folder string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: "image",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
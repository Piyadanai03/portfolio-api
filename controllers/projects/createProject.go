package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils" // 🌟 Import โฟลเดอร์ utils
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProject godoc
// @Summary      เพิ่มข้อมูลโปรเจกต์
// @Description  บันทึกข้อมูลโปรเจกต์ใหม่พร้อมอัปโหลดรูปภาพ (ต้อง Login)
// @Tags         Projects
// @Accept       multipart/form-data
// @Produce      json
// @Param        title        formData string true  "ชื่อโปรเจกต์"
// @Param        description  formData string false "รายละเอียดโปรเจกต์"
// @Param        github_url   formData string false "ลิงก์ Github"
// @Param        cover_image  formData file   false "รูปภาพหน้าปกโปรเจกต์"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/projects [post]
// @Security     BearerAuth
func CreateProject(c *gin.Context) {
	// 1. รับข้อมูลจาก Form-Data
	title := c.PostForm("title")
	description := c.PostForm("description")
	githubURL := c.PostForm("github_url")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุชื่อโปรเจกต์ (title)"})
		return
	}

	// 2. จัดการอัปโหลดรูปภาพหน้าปก (ถ้ามี)
	var coverImageURL string
	file, _, err := c.Request.FormFile("cover_image")
	if err == nil {
		defer file.Close()
		// 🌟 เรียกใช้ฟังก์ชันจาก utils
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_projects")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดรูปหน้าปกล้มเหลว: " + uploadErr.Error()})
			return
		}
		coverImageURL = uploadedURL
	}

	// 3. ดึง User ID จาก Token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลยืนยันตัวตน"})
		return
	}
	uID, _ := uuid.Parse(userID.(string))

	// 4. บันทึกลงฐานข้อมูล
	project := models.Project{
		UserID:        uID,
		Title:         title,
		Description:   description,
		CoverImageURL: coverImageURL,
		GithubURL:     githubURL,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างโปรเจกต์ไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มโปรเจกต์สำเร็จ!", "data": project})
}
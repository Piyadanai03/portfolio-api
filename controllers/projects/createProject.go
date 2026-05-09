package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateProject(c *gin.Context) {
	// รับข้อมูลจาก Form-Data แบบพื้นฐาน
	title := c.PostForm("title")
	description := c.PostForm("description")
	githubURL := c.PostForm("githubURL")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุชื่อโปรเจกต์ (title)"})
		return
	}

	// จัดการอัปโหลด Cover Image
	var coverImageURL string
	file, _, err := c.Request.FormFile("coverImage")
	if err == nil {
		defer file.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_projects")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดรูปหน้าปกล้มเหลว"})
			return
		}
		coverImageURL = uploadedURL
	}

	userID, _ := c.Get("user_id")
	uID, _ := uuid.Parse(userID.(string))

	// 3. สร้าง Object โปรเจกต์
	project := models.Project{
		UserID:        uID,
		Title:         title,
		Description:   description,
		CoverImageURL: coverImageURL,
		GithubURL:     githubURL,
	}

	// ดึงข้อมูล Technologies จาก techIds ที่ส่งมาเป็น Array
	techIDs := c.PostFormArray("techIds")
	if len(techIDs) > 0 {
		var techs []models.Technology
		config.DB.Where("id IN ?", techIDs).Find(&techs)
		project.Technologies = techs // GORM จะจัดการ Many-to-Many ให้เอง
	}

	// บันทึกตัวโปรเจกต์ก่อน เพื่อให้ได้ ProjectID มาใช้กับรูปภาพแกลลอรี
	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างโปรเจกต์ไม่สำเร็จ"})
		return
	}

	// จัดการ Gallery Images (ลูปอัปโหลดหลายไฟล์)
	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["galleryImages"]
		captions := form.Value["galleryCaptions"]

		for i, fileHeader := range files {
			f, err := fileHeader.Open()
			if err != nil {
				continue
			}
			
			uploadedURL, uploadErr := utils.UploadToCloudinary(f, "portfolio_gallery")
			f.Close()
			
			if uploadErr == nil {
				caption := ""
				if i < len(captions) {
					caption = captions[i]
				}
				// บันทึกลงตาราง project_images
				galleryImage := models.ProjectImage{
					ProjectID: project.ID,
					ImageURL:  uploadedURL,
					Caption:   caption,
				}
				config.DB.Create(&galleryImage)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มโปรเจกต์สำเร็จ!", "data": project})
}
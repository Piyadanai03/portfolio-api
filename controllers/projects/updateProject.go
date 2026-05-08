package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils" // 🌟 Import โฟลเดอร์ utils
	"github.com/gin-gonic/gin"
)

// UpdateProject godoc
// @Summary      อัปเดตข้อมูลโปรเจกต์
// @Description  แก้ไขข้อมูลโปรเจกต์ตาม ID (ต้อง Login)
// @Tags         Projects
// @Accept       multipart/form-data
// @Produce      json
// @Param        id           path     string true  "ID ของโปรเจกต์"
// @Param        title        formData string false "ชื่อโปรเจกต์"
// @Param        description  formData string false "รายละเอียดโปรเจกต์"
// @Param        github_url   formData string false "ลิงก์ Github"
// @Param        cover_image  formData file   false "อัปโหลดรูปหน้าปกใหม่ (ถ้าไม่เปลี่ยนไม่ต้องส่ง)"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/projects/{id} [put]
// @Security     BearerAuth
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	// 1. ตรวจสอบว่ามีโปรเจกต์นี้อยู่จริง
	if err := config.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโปรเจกต์ที่ต้องการแก้ไข"})
		return
	}

	// 2. รับข้อมูลใหม่แบบ Form-Data
	title := c.PostForm("title")
	description := c.PostForm("description")
	githubURL := c.PostForm("github_url")

	// 3. จัดการอัปโหลดไฟล์รูปภาพใหม่ (ถ้ามีส่งมา)
	coverImageURL := project.CoverImageURL // ใช้ของเดิมยืนพื้น
	file, _, err := c.Request.FormFile("cover_image")
	if err == nil {
		defer file.Close()
		// 🌟 เรียกใช้ฟังก์ชันจาก utils
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_projects")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดรูปหน้าปกใหม่ล้มเหลว: " + uploadErr.Error()})
			return
		}
		coverImageURL = uploadedURL // เปลี่ยนไปใช้รูปใหม่
	}

	// 4. นำข้อมูลใหม่ไปทับของเดิม
	if title != "" {
		project.Title = title
	}
	if description != "" {
		project.Description = description
	}
	if githubURL != "" {
		project.GithubURL = githubURL
	}
	project.CoverImageURL = coverImageURL

	// 5. บันทึกลงฐานข้อมูล
	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปเดตข้อมูลไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตข้อมูลสำเร็จ!", "data": project})
}
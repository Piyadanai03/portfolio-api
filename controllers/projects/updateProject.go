package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateProject godoc
// @Summary      แก้ไขข้อมูลโปรเจกต์
// @Description  อัพเดตข้อมูลโปรเจกต์ รวมถึงรูปปก แกลลอรี และเทคโนโลยีที่ใช้
// @Tags         Projects
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      string  true  "ID ของโปรเจกต์"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /member/projects/{id} [put]
// @Security     BearerAuth
func UpdateProject(c *gin.Context) {
	id := c.Param("id")

	// ค้นหาโปรเจกต์เดิม พร้อมดึงความสัมพันธ์มาให้ครบ
	var project models.Project
	if err := config.DB.Preload("Technologies").Preload("Images").Preload("Experience").First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโปรเจกต์"})
		return
	}

	// 1. อัปเดตข้อมูล Text พื้นฐาน
	if title := c.PostForm("title"); title != "" {
		project.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		project.Description = description
	}
	project.GithubURL = c.PostForm("githubURL")

	// 2. จัดการ Cover Image
	// ถ้ามีการอัปโหลดไฟล์รูปปกใหม่ ไฟล์จะทับ URL เสมอ
	if file, _, err := c.Request.FormFile("coverImage"); err == nil {
		defer file.Close()
		if uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_projects"); uploadErr == nil {
			project.CoverImageURL = uploadedURL
		}
	} else {
		// ไม่มีไฟล์ใหม่ -> ใช้ค่าจากฟอร์มแทน (ว่าง/null แปลว่าผู้ใช้กดลบรูปทิ้ง)
		coverImageURLInput := c.PostForm("coverImageURL")
		if coverImageURLInput == "" || coverImageURLInput == "null" {
			project.CoverImageURL = ""
		} else {
			project.CoverImageURL = coverImageURLInput
		}
	}

	// 3. จัดการ Experience
	// เคลียร์ object ความสัมพันธ์เดิมที่ Preload มา ป้องกัน GORM sync FK เก่ากลับ
	project.Experience = nil

	experienceID := c.PostForm("experienceID")
	if experienceID != "" && experienceID != "null" && experienceID != "undefined" {
		if parsedExpID, err := uuid.Parse(experienceID); err == nil {
			project.ExperienceID = &parsedExpID
		}
	} else {
		project.ExperienceID = nil
	}

	// บันทึกข้อมูลหลัก (Text, รูปปก, Experience) — Omit("Experience") กัน GORM ยุ่งกับ association
	config.DB.Omit("Experience").Save(&project)

	if project.ExperienceID == nil {
		config.DB.Exec("UPDATE projects SET experience_id = NULL WHERE id = ?", project.ID)
	}

	// 4. จัดการ Technologies (ลบของเก่าแล้วเอาที่เลือกใหม่ใส่แทน)
	if techIDs := c.PostFormArray("techIds"); len(techIDs) > 0 {
		var techs []models.Technology
		config.DB.Where("id IN ?", techIDs).Find(&techs)
		config.DB.Model(&project).Association("Technologies").Replace(techs)
	} else {
		config.DB.Model(&project).Association("Technologies").Clear()
	}

	// 5. จัดการลบรูปภาพ Gallery ที่โดนกดกากบาททิ้ง
	if deletedGalleryIds := c.PostFormArray("deletedGalleryIds"); len(deletedGalleryIds) > 0 {
		config.DB.Where("id IN ?", deletedGalleryIds).Delete(&models.ProjectImage{})
	}

	// 6. จัดการอัปเดตคำบรรยายรูป (Caption) ของรูป Gallery เดิม
	existingImageIds := c.PostFormArray("existingImageIds")
	existingImageCaptions := c.PostFormArray("existingImageCaptions")
	for i, imgID := range existingImageIds {
		if i < len(existingImageCaptions) {
			config.DB.Model(&models.ProjectImage{}).Where("id = ?", imgID).Update("caption", existingImageCaptions[i])
		}
	}

	// 7. จัดการอัปโหลดรูปภาพ Gallery ใหม่ที่เพิ่มเข้ามา
	if form, err := c.MultipartForm(); err == nil {
		newFiles := form.File["galleryImages"]
		newCaptions := form.Value["galleryCaptions"]

		for i, fileHeader := range newFiles {
			f, err := fileHeader.Open()
			if err != nil {
				continue
			}

			uploadedURL, uploadErr := utils.UploadToCloudinary(f, "portfolio_gallery")
			f.Close()

			if uploadErr != nil {
				continue
			}

			caption := ""
			if i < len(newCaptions) {
				caption = newCaptions[i]
			}

			config.DB.Create(&models.ProjectImage{
				ProjectID: project.ID,
				ImageURL:  uploadedURL,
				Caption:   caption,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตข้อมูลสำเร็จ!"})
}
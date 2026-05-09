package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := config.DB.Preload("Technologies").Preload("Images").First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโปรเจกต์"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	githubURL := c.PostForm("githubURL")

	// 1. Cover Image
	coverImageURL := project.CoverImageURL
	file, _, err := c.Request.FormFile("coverImage")
	if err == nil {
		defer file.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_projects")
		if uploadErr == nil {
			coverImageURL = uploadedURL
		}
	}

	// อัปเดตข้อมูล Text พื้นฐาน
	if title != "" { project.Title = title }
	if description != "" { project.Description = description }
	project.GithubURL = githubURL
	project.CoverImageURL = coverImageURL
	config.DB.Save(&project)

	techIDs := c.PostFormArray("techIds")
	var techs []models.Technology
	if len(techIDs) > 0 {
		config.DB.Where("id IN ?", techIDs).Find(&techs)
	}
	config.DB.Model(&project).Association("Technologies").Replace(techs)

	deletedGalleryIds := c.PostFormArray("deletedGalleryIds")
	if len(deletedGalleryIds) > 0 {
		config.DB.Where("id IN ?", deletedGalleryIds).Delete(&models.ProjectImage{})
	}

	existingImageIds := c.PostFormArray("existingImageIds")
	existingImageCaptions := c.PostFormArray("existingImageCaptions")
	for i, imgID := range existingImageIds {
		if i < len(existingImageCaptions) {
			config.DB.Model(&models.ProjectImage{}).Where("id = ?", imgID).Update("caption", existingImageCaptions[i])
		}
	}

	form, err := c.MultipartForm()
	if err == nil {
		newFiles := form.File["galleryImages"]
		newCaptions := form.Value["galleryCaptions"]

		for i, fileHeader := range newFiles {
			f, err := fileHeader.Open()
			if err != nil { continue }

			uploadedURL, uploadErr := utils.UploadToCloudinary(f, "portfolio_gallery")
			f.Close()

			if uploadErr == nil {
				caption := ""
				if i < len(newCaptions) { caption = newCaptions[i] }
				
				config.DB.Create(&models.ProjectImage{
					ProjectID: project.ID,
					ImageURL:  uploadedURL,
					Caption:   caption,
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตข้อมูลสำเร็จ!"})
}
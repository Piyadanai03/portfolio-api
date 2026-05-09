package projects

import (
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	result := config.DB.Preload("Images").Preload("Technologies").Order("created_at desc").Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func GetProjectByID(c *gin.Context) {
	var project models.Project

	id := c.Param("id")

	result := config.DB.
		Preload("Images").
		Preload("Technologies").
		First(&project, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบโปรเจค",
		})
		return
	}

	c.JSON(http.StatusOK, project)
}

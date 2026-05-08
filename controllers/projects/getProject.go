package projects

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	// 🌟 เติม .Preload("Technologies") เพื่อดึง Tech Stack มาโชว์หน้าตารางด้วย
	result := config.DB.Preload("Images").Preload("Technologies").Order("created_at desc").Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
package projects

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
)

// GetProjects godoc
// @Summary      ดึงข้อมูลโปรเจกต์ทั้งหมด
// @Description  แสดงรายชื่อโปรเจกต์พร้อมรูปภาพประกอบสำหรับหน้าเว็บทั่วไป
// @Tags         Projects
// @Produce      json
// @Success      200  {array}   models.Project
// @Router       /projects [get]
func GetProjects(c *gin.Context) {
	var projects []models.Project

	result := config.DB.Preload("Images").Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
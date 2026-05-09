package technologies

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

func GetTechnologies(c *gin.Context) {
	var techs []models.Technology

	if err := config.DB.Find(&techs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ดึงข้อมูล technologies ไม่สำเร็จ",
		})
		return
	}

	c.JSON(http.StatusOK, techs)
}
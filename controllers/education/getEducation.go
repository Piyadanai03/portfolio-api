package education

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// GetEducations godoc
// @Summary      ดึงรายการการศึกษา
// @Description  ดึงข้อมูลการศึกษาทั้งหมด
// @Tags         Education
// @Produce      json
// @Success      200  {array}   models.Study
// @Failure      500  {object}  map[string]interface{}
// @Router       /educations [get]
func GetEducations(c *gin.Context) {
	var studies []models.Study

	result := config.DB.Order("graduation_date desc").Find(&studies)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลได้"})
		return
	}

	c.JSON(http.StatusOK, studies)
}

// GetEducationByID godoc
// @Summary      ดูข้อมูลการศึกษาตาม ID
// @Description  ดึงข้อมูลการศึกษาตาม ID
// @Tags         Education
// @Produce      json
// @Param        id   path      string  true  "ID ของการศึกษา"
// @Success      200  {object}  models.Study
// @Failure      404  {object}  map[string]interface{}
// @Router       /educations/{id} [get]
func GetEducationByID(c *gin.Context) {
	id := c.Param("id")
	var study models.Study

	result := config.DB.First(&study, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลการศึกษา"})
		return
	}

	c.JSON(http.StatusOK, study)
}

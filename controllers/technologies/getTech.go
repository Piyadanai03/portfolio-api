package technologies

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// GetTechnologies godoc
// @Summary      ดึงข้อมูลเทคโนโลยีทั้งหมด
// @Description  ดึงข้อมูลเทคโนโลยีทั้งหมดจากฐานข้อมูล
// @Tags         Technology
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Technology
// @Router       /member/tech [get]
// @Security     BearerAuth
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

// GetTechnologyByID godoc
// @Summary      ดึงข้อมูลเทคโนโลยีตาม ID
// @Description  ดึงข้อมูลเทคโนโลยีจากฐานข้อมูลตาม ID ที่ระบุ
// @Tags         Technology
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID ของเทคโนโลยี"
// @Success      200  {object}  models.Technology
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /member/tech/{id} [get]
// @Security     BearerAuth
func GetTechnologyByID(c *gin.Context) {
	var tech models.Technology
	id := c.Param("id")

	if err := config.DB.First(&tech, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบ technology ที่ระบุ",
		})
		return
	}

	c.JSON(http.StatusOK, tech)
}
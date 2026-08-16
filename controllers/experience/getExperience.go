package experience

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

// GetExperiences godoc
// @Summary      ดึงรายการประสบการณ์ทำงาน
// @Description  ดึงข้อมูลประสบการณ์ทั้งหมด พร้อมโปรเจกต์ที่เกี่ยวข้อง
// @Tags         Experience
// @Produce      json
// @Success      200  {array}   models.Experience
// @Failure      500  {object}  map[string]interface{}
// @Router       /member/experiences [get]
func GetExperiences(c *gin.Context) {
	var experiences []models.Experience

	result := config.DB.Preload("Projects").Order("start_date desc").Find(&experiences)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.Error("not found experiences"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(experiences, "Experiences fetched successfully"))
}

// GetExperienceByID godoc
// @Summary      ดูประสบการณ์ทำงานตาม ID
// @Description  ดึงข้อมูลประสบการณ์ตาม ID
// @Tags         Experience
// @Produce      json
// @Param        id   path      string  true  "ID ของประสบการณ์"
// @Success      200  {object}  models.Experience
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/experiences/{id} [get]
func GetExperienceByID(c *gin.Context) {
	id := c.Param("id")
	var experience models.Experience

	result := config.DB.Preload("Projects").First(&experience, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, utils.Error("not found experience"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(experience, "Experience fetched successfully"))
}

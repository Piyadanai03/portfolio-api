package experience

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
)

// DeleteExperience godoc
// @Summary      ลบข้อมูลประสบการณ์ทำงาน
// @Description  ลบข้อมูลประสบการณ์ทำงานตาม ID (ต้อง Login)
// @Tags         Experience
// @Param        id   path      string  true  "ID ของประสบการณ์ทำงาน"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/experience/{id} [delete]
// @Security     BearerAuth
func DeleteExperience(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Experience{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลสำเร็จ"})
}
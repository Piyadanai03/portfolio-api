package education

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
)

// DeleteEducation godoc
// @Summary      ลบข้อมูลการศึกษา
// @Description  ลบข้อมูลการศึกษาตาม ID (ต้อง Login)
// @Tags         Education
// @Param        id   path      string  true  "ID ของการศึกษา"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/education/{id} [delete]
// @Security     BearerAuth
func DeleteEducation(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Study{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลสำเร็จ"})
}
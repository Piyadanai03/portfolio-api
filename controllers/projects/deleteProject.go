package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// DeleteProject godoc
// @Summary      ลบข้อมูลโปรเจกต์
// @Description  ลบข้อมูลโปรเจกต์ตาม ID (ต้อง Login)
// @Tags         Projects
// @Param        id   path      string  true  "ID ของโปรเจกต์"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/projects/{id} [delete]
// @Security     BearerAuth
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if err := config.DB.First(&project, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลที่ต้องการลบ"})
		return
	}

	config.DB.Model(&project).Association("Technologies").Clear()
	config.DB.Where("project_id = ?", id).Delete(&models.ProjectImage{})
	config.DB.Model(&models.Achievement{}).Where("project_id = ?", id).Update("project_id", nil)

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ลบโปรเจกต์ไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบโปรเจกต์เรียบร้อยแล้ว"})
}
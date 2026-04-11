package projects

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
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
    
    result := config.DB.Delete(&models.Project{}, "id = ?", id)

    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลที่ต้องการลบ"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "ลบโปรเจกต์เรียบร้อยแล้ว"})
}
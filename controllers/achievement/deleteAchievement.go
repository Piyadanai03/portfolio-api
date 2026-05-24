package achievement

import (
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteAchievement godoc
// @Summary      ลบข้อมูลความสำเร็จ
// @Description  ลบข้อมูลความสำเร็จออกจากระบบตาม ID (ต้อง Login)
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID ของความสำเร็จ"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /member/achievement/{id} [delete]
// @Security     BearerAuth
func DeleteAchievement(c *gin.Context) {
	// 1. รับ ID จาก URL Parameter
	id := c.Param("id")

	// 2. ค้นหาข้อมูลในฐานข้อมูลก่อนว่ามีอยู่จริงไหม
	var achievement models.Achievement
	if err := config.DB.First(&achievement, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลความสำเร็จที่ต้องการลบ"})
		return
	}

	// 3. สั่งลบข้อมูลออกจากฐานข้อมูล
	if err := config.DB.Delete(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบข้อมูลความสำเร็จได้"})
		return
	}

	// 4. ส่งข้อความยืนยันการลบ
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลความสำเร็จสำเร็จ!"})
}

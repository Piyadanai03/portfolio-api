package technologies

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// DeleteTech godoc
// @Summary      ลบข้อมูลเทคโนโลยี
// @Description  ลบข้อมูลเทคโนโลยีออกจากระบบตาม ID (ต้อง Login)
// @Tags         Technology
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID ของเทคโนโลยี"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /member/tech/{id} [delete]
// @Security     BearerAuth
func DeleteTech(c *gin.Context) {
	// 1. รับ ID จาก URL Parameter
	id := c.Param("id")

	// 2. ค้นหาข้อมูลในฐานข้อมูลก่อนว่ามีอยู่จริงไหม
	var tech models.Technology
	if err := config.DB.First(&tech, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลเทคโนโลยีที่ต้องการลบ"})
		return
	}

	// 3. สั่งลบข้อมูลออกจากฐานข้อมูล
	if err := config.DB.Delete(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบข้อมูลเทคโนโลยีได้"})
		return
	}

	// 4. ส่งข้อความยืนยันการลบ
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลเทคโนโลยีสำเร็จ!"})
}
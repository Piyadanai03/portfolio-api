package achievement

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// GetAchievement godoc
// @Summary      ดึงข้อมูลความสำเร็จทั้งหมด
// @Description  ดึงข้อมูลความสำเร็จทั้งหมดจากฐานข้อมูล (เรียงจากใหม่ไปเก่า)
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Achievement
// @Router       /member/achievement [get]
// @Security     BearerAuth
func GetAchievement(c *gin.Context) {
	var achievements []models.Achievement

	// 🌟 เพิ่ม .Order("date_achieved desc") เพื่อให้ใบประกาศใหม่ล่าสุดขึ้นก่อน
	if err := config.DB.Order("date_achieved desc").Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ดึงข้อมูล achievement ไม่สำเร็จ",
		})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// GetAchievementByID godoc
// @Summary      ดึงข้อมูลความสำเร็จตาม ID
// @Description  ดึงข้อมูลความสำเร็จตาม ID จากฐานข้อมูล
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Achievement ID" // 🌟 แก้จาก int เป็น string
// @Success      200  {object}  models.Achievement
// @Router       /member/achievement/{id} [get]
// @Security     BearerAuth
func GetAchievementByID(c *gin.Context) {
	var achievement models.Achievement
	id := c.Param("id")

	// 🌟 ปรับ Syntax การค้นหาให้เป็นมาตรฐานเดียวกับไฟล์อื่น ("id = ?", id)
	if err := config.DB.First(&achievement, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบ achievement ที่ระบุ",
		})
		return
	}

	c.JSON(http.StatusOK, achievement)
}
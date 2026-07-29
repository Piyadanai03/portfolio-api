package achievement

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateAchievementRequest struct {
	Title        *string    `json:"title"`
	Category     *string    `json:"category"`
	DateAchieved *time.Time `json:"dateAchieved"`
	ProjectID    *uuid.UUID `json:"projectID"`
}

// UpdateAchievement godoc
// @Summary      แก้ไขข้อมูลความสำเร็จ
// @Description  อัพเดตข้อมูลความสำเร็จ (ต้อง Login)
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "ID ของความสำเร็จ"
// @Param        body  body      UpdateAchievementRequest  true  "Achievement data to update"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /member/achievement/{id} [put]
// @Security     BearerAuth
func UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAchievementRequest

	// Validate request body (เพิ่ม err.Error() เพื่อให้รู้ว่าพังเพราะวันที่รูปแบบผิดหรือเปล่า)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง หรือรูปแบบวันที่ผิด: " + err.Error()})
		return
	}

	// Find existing achievement
	var achievement models.Achievement
	if err := config.DB.First(&achievement, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลความสำเร็จ"})
		return
	}

	// Update fields (Smart Update ด้วย Pointer)
	updateData := map[string]interface{}{}
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Category != nil {
		updateData["category"] = *req.Category
	}
	if req.DateAchieved != nil {
		updateData["date_achieved"] = *req.DateAchieved
	}
	if req.ProjectID != nil {
		updateData["project_id"] = *req.ProjectID
	}

	// Save updates
	if err := config.DB.Model(&achievement).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถแก้ไขความสำเร็จได้"})
		return
	}

	// 🌟 เพิ่มบรรทัดนี้: ดึงข้อมูลที่อัปเดตล่าสุดจาก Database กลับมาใส่ตัวแปร
	config.DB.First(&achievement, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{"message": "แก้ไขความสำเร็จสำเร็จ!", "achievement": achievement})
}
package achievement

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateAchievementRequest struct {
	Title        string     `json:"title" binding:"required"`
	Category     string     `json:"category" binding:"required"`
	DateAchieved time.Time  `json:"dateAchieved" binding:"required"`
	ProjectID    *uuid.UUID `json:"projectID"`
}

// CreateAchievement godoc
// @Summary      สร้างข้อมูลความสำเร็จ
// @Description  สร้างข้อมูลความสำเร็จใหม่ (ต้อง Login)
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Param        body  body      CreateAchievementRequest  true  "Achievement data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /member/achievement [post]
// @Security     BearerAuth
func CreateAchievement(c *gin.Context) {
	var req CreateAchievementRequest

	// 1. Validate request body (เพิ่ม err.Error() เพื่อให้รู้ว่าพังตรงไหน)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน หรือรูปแบบวันที่ผิด (ต้องเป็น ISO8601): " + err.Error()})
		return
	}

	// 🌟 2. แก้ไขการดึง UserID ให้เหมือนกับ Module อื่นๆ
	userIDStr, exists := c.Get("user_id") // ใช้ user_id ตาม Middleware ที่เราตั้งไว้
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบข้อมูลผู้ใช้"})
		return
	}

	uID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ไม่ถูกต้อง"})
		return
	}

	// 3. Create achievement
	achievement := models.Achievement{
		Title:        req.Title,
		Category:     req.Category,
		DateAchieved: req.DateAchieved,
		ProjectID:    req.ProjectID,
		UserID:       uID, // ใช้ uID ที่แปลงเป็น UUID เรียบร้อยแล้ว
	}

	// 4. Save to database
	if err := config.DB.Create(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างข้อมูลความสำเร็จไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "สร้างข้อมูลความสำเร็จสำเร็จ!", "achievement": achievement})
}

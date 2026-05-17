package experience

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateExperience godoc
// @Summary      เพิ่มข้อมูลประสบการณ์ทำงาน
// @Description  บันทึกข้อมูลประสบการณ์ทำงานใหม่ (ต้อง Login)
// @Tags         Experience
// @Accept       json
// @Produce      json
// @Param        input  body  object  true  "ข้อมูลประสบการณ์ทำงาน"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/experience [post]
// @Security     BearerAuth
func CreateExperience(c *gin.Context) {
	var input struct {
		JobTitle    string     `json:"jobTitle" binding:"required"`
		Company     string     `json:"company" binding:"required"`
		StartDate   time.Time  `json:"startDate" binding:"required"`
		EndDate     *time.Time `json:"endDate"` // ใช้ Pointer รับค่า null ได้ (ถ้ายังทำงานอยู่)
		Description string     `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน หรือรูปแบบวันที่ไม่ถูกต้อง (ต้องเป็น ISO 8601): " + err.Error()})
		return
	}

	// 2. ดึง User ID จาก Token
	userID, _ := c.Get("user_id")
	uID, _ := uuid.Parse(userID.(string))

	// 🌟 3. นำข้อมูลใส่ Model จริง (ครอบคลุมทั้งกรณีมี/ไม่มี EndDate)
	experience := models.Experience{
		UserID:      uID,
		JobTitle:    input.JobTitle,
		Company:     input.Company,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate, 
		Description: input.Description,
	}

	// 4. บันทึกลงฐานข้อมูล
	if err := config.DB.Create(&experience).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มข้อมูลประสบการณ์ทำงานสำเร็จ!", "data": experience})
}
package experience

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// UpdateExperience godoc
// @Summary      แก้ไขข้อมูลประสบการณ์ทำงาน
// @Description  อัพเดตข้อมูลประสบการณ์ทำงาน (ต้อง Login)
// @Tags         Experience
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "ID ของประสบการณ์ทำงาน"
// @Param        input  body      object  true  "ข้อมูลประสบการณ์ทำงานที่ต้องการแก้ไข"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /member/experience/{id} [put]
// @Security     BearerAuth
func UpdateExperience(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		JobTitle    string     `json:"jobTitle"`
		Company     string     `json:"company"`
		StartDate   time.Time  `json:"startDate"`
		EndDate     *time.Time `json:"endDate"` // ใช้ Pointer รับค่า null ได้ (ถ้ายังทำงานอยู่)
		Description string     `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง หรือรูปแบบวันที่ไม่ถูกต้อง (ต้องเป็น ISO 8601): " + err.Error()})
		return
	}

	// 1. ค้นหาข้อมูลประสบการณ์ทำงานเดิม
	var experience models.Experience
	if err := config.DB.First(&experience, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลประสบการณ์ทำงาน"})
		return
	}

	// 🌟 2. Smart Update: เช็คค่าก่อนจับใส่ Map
	updateData := map[string]interface{}{}

	if input.JobTitle != "" {
		updateData["job_title"] = input.JobTitle
	}
	if input.Company != "" {
		updateData["company"] = input.Company
	}
	if input.Description != "" {
		updateData["description"] = input.Description
	}
	if !input.StartDate.IsZero() {
		updateData["start_date"] = input.StartDate
	}
	
	updateData["end_date"] = input.EndDate

	if err := config.DB.Model(&experience).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัพเดตข้อมูลไม่สำเร็จ"})
		return
	}

	config.DB.First(&experience, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{"message": "แก้ไขข้อมูลประสบการณ์ทำงานสำเร็จ!", "data": experience})
}
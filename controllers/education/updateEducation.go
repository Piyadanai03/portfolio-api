package education

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
)

// UpdateEducation godoc
// @Summary      แก้ไขข้อมูลการศึกษา
// @Description  อัพเดตข้อมูลการศึกษาตาม ID (ต้อง Login)
// @Tags         Education
// @Accept       json
// @Produce      json
// @Param        id     path      string  true  "ID ของการศึกษา"
// @Param        input  body      object  true  "ข้อมูลการศึกษาที่ต้องการแก้ไข"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /member/education/{id} [put]
// @Security     BearerAuth
func UpdateEducation(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Degree         string    `json:"degree"`
		Faculty        string    `json:"faculty"`
		Major          string    `json:"major"`
		Institution    string    `json:"institution"`
		GPA            *float64  `json:"gpa"`
		GraduationDate time.Time `json:"graduationDate"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง หรือรูปแบบไม่ถูกต้อง: " + err.Error()})
		return
	}

	var study models.Study
	if err := config.DB.First(&study, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลการศึกษา"})
		return
	}

	updateData := map[string]interface{}{}

	if input.Degree != "" {
		updateData["degree"] = input.Degree
	}
	if input.Faculty != "" {
		updateData["faculty"] = input.Faculty
	}
	if input.Major != "" {
		updateData["major"] = input.Major
	}
	if input.Institution != "" {
		updateData["institution"] = input.Institution
	}
	
	updateData["gpa"] = input.GPA

	if !input.GraduationDate.IsZero() {
		updateData["graduation_date"] = input.GraduationDate
	}

	if err := config.DB.Model(&study).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัพเดตข้อมูลไม่สำเร็จ"})
		return
	}

	config.DB.First(&study, "id = ?", id) 

	c.JSON(http.StatusOK, gin.H{"message": "แก้ไขข้อมูลการศึกษาสำเร็จ!", "data": study})
}

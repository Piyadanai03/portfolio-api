package education

import (
	"net/http"
	"time"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEducation godoc
// @Summary      เพิ่มข้อมูลการศึกษา
// @Description  บันทึกข้อมูลการศึกษาใหม่ (ต้อง Login)
// @Tags         Education
// @Accept       json
// @Produce      json
// @Param        input  body  models.Study  true  "ข้อมูลการศึกษา"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/education [post]
// @Security     BearerAuth
func CreateEducation(c *gin.Context) {
	var input struct {
		Degree         string    `json:"degree" binding:"required"`
		Faculty        string    `json:"faculty"`
		Major          string    `json:"major" binding:"required"`
		Institution    string    `json:"institution" binding:"required"`
		GPA            *float64  `json:"gpa"`
		GraduationDate time.Time `json:"graduationDate"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ครบถ้วน หรือรูปแบบไม่ถูกต้อง: " + err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uID, _ := uuid.Parse(userID.(string))

	study := models.Study{
		UserID:         uID,
		Degree:         input.Degree,
		Faculty:        input.Faculty,
		Major:          input.Major,
		Institution:    input.Institution,
		GPA:            input.GPA,
		GraduationDate: input.GraduationDate,
	}

	if err := config.DB.Create(&study).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลการศึกษาไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มข้อมูลการศึกษาสำเร็จ!", "data": study})
}
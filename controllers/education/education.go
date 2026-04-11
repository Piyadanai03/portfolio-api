package education

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/google/uuid"
)

// CreateEducation godoc
// @Summary      เพิ่มข้อมูลการศึกษา
// @Description  บันทึกข้อมูลการศึกษาใหม่ (ต้อง Login)
// @Tags         Education
// @Accept       json
// @Produce      json
// @Param        input  body  models.Study  true  "ข้อมูลการศึกษา"
// @Success      201    {object}  models.Study
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/education [post]
// @Security     BearerAuth
func CreateEducation(c *gin.Context) {
	var input models.Study
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uID, _ := uuid.Parse(userID.(string))
	input.UserID = uID

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกข้อมูลไม่สำเร็จ"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// DeleteEducation godoc
// @Summary      ลบข้อมูลการศึกษา
// @Description  ลบข้อมูลการศึกษาตาม ID (ต้อง Login)
// @Tags         Education
// @Param        id   path      string  true  "ID ของการศึกษา"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/education/{id} [delete]
// @Security     BearerAuth
func DeleteEducation(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Study{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลสำเร็จ"})
}
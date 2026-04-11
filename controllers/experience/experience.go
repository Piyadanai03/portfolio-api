package experience

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/google/uuid"
)

// CreateExperience godoc
// @Summary      เพิ่มข้อมูลประสบการณ์ทำงาน
// @Description  บันทึกข้อมูลประสบการณ์ทำงานใหม่ (ต้อง Login)
// @Tags         Experience
// @Accept       json
// @Produce      json
// @Param        input  body  models.Experience  true  "ข้อมูลประสบการณ์ทำงาน"
// @Success      201    {object}  models.Experience
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/experience [post]
// @Security     BearerAuth
func CreateExperience(c *gin.Context) {
	var input models.Experience
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

// DeleteExperience godoc
// @Summary      ลบข้อมูลประสบการณ์ทำงาน
// @Description  ลบข้อมูลประสบการณ์ทำงานตาม ID (ต้อง Login)
// @Tags         Experience
// @Param        id   path      string  true  "ID ของประสบการณ์ทำงาน"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/experience/{id} [delete]
// @Security     BearerAuth
func DeleteExperience(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Experience{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลสำเร็จ"})
}
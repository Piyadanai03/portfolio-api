package projects

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/google/uuid"
)

// CreateProject godoc
// @Summary      เพิ่มข้อมูลโปรเจกต์
// @Description  บันทึกข้อมูลโปรเจกต์ใหม่ (ต้อง Login)
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        input  body  object  true  "ข้อมูลโปรเจกต์"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /member/projects [post]
// @Security     BearerAuth
func CreateProject(c *gin.Context) {
	var input struct {
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description"`
		CoverImageURL string `json:"cover_image_url"`
		GithubURL     string `json:"github_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ดึง User ID จาก Middleware (ที่ยามฝากไว้ใน Context)
	userID, _ := c.Get("user_id")
	
	// แปลง Interface เป็น uuid.UUID
	uID, _ := uuid.Parse(userID.(string))

	project := models.Project{
		UserID:        uID,
		Title:         input.Title,
		Description:   input.Description,
		CoverImageURL: input.CoverImageURL,
		GithubURL:     input.GithubURL,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างโปรเจกต์ไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "เพิ่มโปรเจกต์สำเร็จ!", "data": project})
}
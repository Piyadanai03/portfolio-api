package projects

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
)

// UpdateProject godoc
// @Summary      อัปเดตข้อมูลโปรเจกต์
// @Description  แก้ไขข้อมูลโปรเจกต์ตาม ID (ต้อง Login)
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID ของโปรเจกต์"
// @Param        input  body  object  true  "ข้อมูลโปรเจกต์ที่ต้องการแก้ไข"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /member/projects/{id} [put]
// @Security     BearerAuth
func UpdateProject(c *gin.Context) {
    id := c.Param("id") // รับ ID จาก URL เช่น /projects/1
    var project models.Project

    // 1. หาโปรเจกต์ใน DB ก่อนว่ามีอยู่จริงไหม
    if err := config.DB.First(&project, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบโปรเจกต์ที่ต้องการแก้ไข"})
        return
    }

    // 2. รับข้อมูลใหม่จาก Body
    var input struct {
        Title         string `json:"title"`
        Description   string `json:"description"`
        CoverImageURL string `json:"cover_image_url"`
        GithubURL     string `json:"github_url"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 3. อัปเดตเฉพาะฟิลด์ที่ส่งมา
    config.DB.Model(&project).Updates(input)

    c.JSON(http.StatusOK, gin.H{"message": "อัปเดตข้อมูลสำเร็จ!", "data": project})
}
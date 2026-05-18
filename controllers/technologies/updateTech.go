package technologies

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

// UpdateTech godoc
// @Summary      แก้ไขข้อมูลเทคโนโลยี
// @Description  อัพเดตข้อมูลเทคโนโลยี (เช่น ชื่อ, หมวดหมู่, และเปลี่ยนรูปไอคอน) (ต้อง Login)
// @Tags         Technology
// @Accept       multipart/form-data
// @Produce      json
// @Param        id        path      string  true  "ID ของเทคโนโลยี"
// @Param        name      formData  string  false "ชื่อเทคโนโลยี"
// @Param        category  formData  string  false "หมวดหมู่เทคโนโลยี"
// @Param        icon      formData  file    false "ไฟล์รูปภาพไอคอน"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Router       /member/tech/{id} [put]
// @Security     BearerAuth
func UpdateTech(c *gin.Context) {
	// 🌟 1. รับ ID จาก URL (เป็น String ได้เลย ไม่ต้องแปลงเป็นตัวเลข)
	id := c.Param("id")

	// 2. ค้นหาเทคโนโลยีในฐานข้อมูล (ใช้ String ค้นหา UUID แบบนี้ได้เลย)
	var tech models.Technology
	if err := config.DB.First(&tech, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลเทคโนโลยี"})
		return
	}

	// 3. รับข้อมูลจาก Form-Data แบบพื้นฐาน
	name := c.PostForm("name")
	category := c.PostForm("category")

	// อัปเดตข้อมูลเทคโนโลยี (ถ้ามีการส่งค่ามา)
	if name != "" {
		tech.Name = name
	}
	if category != "" {
		tech.Category = category
	}

	// 4. จัดการอัปโหลด Icon Image
	file, _, err := c.Request.FormFile("icon")
	if err == nil {
		defer file.Close()
		uploadedURL, uploadErr := utils.UploadToCloudinary(file, "portfolio_technologies")
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปโหลดภาพได้"})
			return
		}
		tech.IconURL = uploadedURL // ทับ URL เดิมด้วยรูปใหม่
	}

	// 5. บันทึกการอัปเดตลงในฐานข้อมูล
	if err := config.DB.Save(&tech).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตข้อมูลเทคโนโลยีได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "อัปเดตเทคโนโลยีสำเร็จ!", "technology": tech})
}
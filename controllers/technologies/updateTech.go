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
	// 1. รับ ID จาก URL
	id := c.Param("id")

	// 2. ค้นหาเทคโนโลยีในฐานข้อมูล
	var tech models.Technology
	if err := config.DB.First(&tech, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลเทคโนโลยี"})
		return
	}

	// 3. รับข้อมูลจาก Form-Data
	name := c.PostForm("name")
	category := c.PostForm("category")
	// 🌟 เพิ่มบรรทัดนี้: รับค่า iconURL ด้วย
	iconURL := c.PostForm("iconURL")

	// อัปเดตข้อมูลเทคโนโลยี (ถ้ามีการส่งค่ามา)
	if name != "" {
		tech.Name = name
	}
	if category != "" {
		tech.Category = category
	}
	// 🌟 เพิ่มบรรทัดนี้: อัปเดต iconURL ถ้ามีการส่งข้อความ URL มา
	if iconURL != "" {
		tech.IconURL = iconURL
	}

	// 4. จัดการอัปโหลด Icon Image (ถ้ามีการแนบไฟล์มา จะทับ iconURL ข้างบนเสมอ)
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
package projects

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// 1. ดึงไฟล์จาก Form-data
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่พบไฟล์รูปภาพ"})
		return
	}
	defer file.Close()

	// 2. ตั้งค่าเชื่อมต่อกับ Cloudinary (ดึงค่าจาก .env)
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เชื่อมต่อ Cloudinary ไม่สำเร็จ"})
		return
	}

	// 3. กำหนด Context และส่งไฟล์ขึ้น Cloud
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "portfolio_projects", // ชื่อโฟลเดอร์ใน Cloudinary
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปโหลดล้มเหลว: " + err.Error()})
		return
	}

	// 4. ส่ง URL ที่ได้จาก Cloud กลับไปให้ผู้ใช้
	c.JSON(http.StatusOK, gin.H{
		"message":   "อัปโหลดสำเร็จ!",
		"image_url": uploadResult.SecureURL, // นี่คือ URL https ที่เราต้องการ
	})
}
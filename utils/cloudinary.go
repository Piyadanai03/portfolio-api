package utils

import (
	"context"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadToCloudinary อัปโหลดไฟล์ไปยัง Cloudinary (เปิดเป็น Public ให้โฟลเดอร์อื่นเรียกใช้ได้)
func UploadToCloudinary(file interface{}, folder string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	// กำหนดเวลา Timeout ป้องกันเน็ตหลุดหรือประมวลผลนานเกินไป (30 วินาที)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:       folder,
		ResourceType: "image", // ใช้ "image" ได้เลยเพราะรับได้ทั้งภาพและ PDF ที่ตั้งค่าไว้
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
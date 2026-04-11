package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"os"
)

// Login godoc
// @Summary      เข้าสู่ระบบ (Admin)
// @Description  ตรวจสอบ Username/Password และคืนค่า JWT Token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        input  body  object  true  "ข้อมูล Login (username, password)"
// @Success      200    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /login [post]
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณากรอกข้อมูลให้ครบ"})
		return
	}

	// 1. ค้นหา User ใน DB
	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบผู้ใช้นี้ หรือรหัสผ่านไม่ถูกต้อง"})
		return
	}

	// 2. ตรวจสอบรหัสผ่านที่ Hash ไว้
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านไม่ถูกต้อง"})
		return
	}

	// 3. สร้าง JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // บัตรมีอายุ 24 ชม.
	})

	// เซ็นชื่อกำกับบัตรด้วยความลับ (Secret Key)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(http.StatusOK, gin.H{
		"message": "Login สำเร็จ!",
		"token":   tokenString,
	})
}
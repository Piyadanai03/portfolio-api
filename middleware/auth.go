package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Token จาก Header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "กรุณา Login ก่อนใช้งาน"})
			c.Abort() // หยุดการทำงานทันที ไม่ให้ไปต่อที่ Controller
			return
		}

		// รูปแบบปกติคือ "Bearer [TOKEN]" เราต้องตัดคำว่า Bearer ออก
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. ตรวจสอบความถูกต้องของ Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ไม่ถูกต้องหรือหมดอายุ"})
			c.Abort()
			return
		}

		// 3. ถ้าผ่าน ให้เอา user_id ไปใส่ใน Context เผื่อ Controller ต้องใช้
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])

		c.Next() // อนุญาตให้ผ่านไปที่ Controller ได้
	}
}
package portfolio

import (
	"net/http"
	"os"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
	
)

func GetHomeData(c *gin.Context) {
	var user models.User

	ownerID := os.Getenv("USER_ID")

	if ownerID == "" {
		c.JSON(http.StatusInternalServerError, utils.Error("not found ownerID"))
		return
	}

	if err := config.DB.Where("id = ?", ownerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error("user not found"))
		return
	}

	// 🌟 เช็คให้แน่ใจว่ามีเงื่อนไข WHERE is_featured = true และ Limit(3) ตรงนี้
	if err := config.DB.Where("user_id = ? AND is_featured = true", ownerID).Limit(3).Find(&user.Projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error("failed to fetch featured projects"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}
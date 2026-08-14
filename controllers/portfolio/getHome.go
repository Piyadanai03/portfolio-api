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

	if err := config.DB.Preload("Projects").Where("id = ?", ownerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error("user not found"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}
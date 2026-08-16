package portfolio

import (
	"net/http"
	"os"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

// GetAbout godoc
// @Summary Get about data
// @Description Get about data for the portfolio
// @Tags Portfolio
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/about [get]
func GetAbout(c *gin.Context) {
	var user models.User
	var achievements []models.Achievement

	ownerID := os.Getenv("USER_ID")

	if ownerID == "" {
		c.JSON(http.StatusInternalServerError, utils.Error("not found ownerID"))
		return
	}

	if err := config.DB.
		Preload("Experiences").
		Preload("Studies").
		Preload("Contacts").
		Preload("Projects").
		Where("id = ?", ownerID).
		First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.Error("not found user data"))
		return
	}

	if err := config.DB.
		Where("user_id = ?", ownerID).
		Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.Error("not found achievements"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(gin.H{
		"user":         user,
		"achievements": achievements,
	}, "About data fetched successfully"))
}

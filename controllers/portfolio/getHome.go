package portfolio

import (
	"net/http"
	"os"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

// GetHomeData godoc
// @Summary Get home data
// @Description Get home data for the portfolio
// @Tags Portfolio
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/home [get]
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

	if err := config.DB.
		Preload("Images").
		Preload("Technologies").
		Preload("Experience").
		Preload("Achievements").
		Where("user_id = ? AND is_featured = true", ownerID).
		Limit(3).
		Find(&user.Projects).Error; err != nil {

		c.JSON(http.StatusInternalServerError, utils.Error("failed to fetch featured projects"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user, "Home data fetched successfully"))
}

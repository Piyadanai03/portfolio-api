package projects

import (
	"net/http"

	"github.com/Piyadanai03/portfolio-api/config"
	"github.com/Piyadanai03/portfolio-api/models"
	"github.com/Piyadanai03/portfolio-api/utils"
	"github.com/gin-gonic/gin"
)

func GetProjects(c *gin.Context) {
	var projects []models.Project

	result := config.DB.Preload("Images").Preload("Technologies").Preload("Experience").Preload("Achievements").Order("created_at desc").Find(&projects)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.Error("failed to fetch projects"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(projects, "projects fetched successfully"))
}

func GetProjectByID(c *gin.Context) {
	var project models.Project
	id := c.Param("id")

	result := config.DB.
		Preload("Images").
		Preload("Technologies").
		Preload("Experience").
		Preload("Achievements").
		Where("id = ?", id).
		First(&project)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, utils.Error("project not found"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(project, "project fetched successfully"))
}

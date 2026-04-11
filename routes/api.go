package routes

import (
	_ "github.com/Piyadanai03/portfolio-api/docs"
	"github.com/gin-gonic/gin"
	"github.com/Piyadanai03/portfolio-api/controllers/auth"
	"github.com/Piyadanai03/portfolio-api/controllers/projects"
	"github.com/Piyadanai03/portfolio-api/middleware"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/Piyadanai03/portfolio-api/controllers/education"
	"github.com/Piyadanai03/portfolio-api/controllers/experience"
	
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api/v1")

	v1.GET("/projects", projects.GetProjects)
	v1.POST("/login", auth.Login)

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		
	}

	member := v1.Group("/member")
	member.Use(middleware.AuthMiddleware())
	{
		member.POST("/projects", projects.CreateProject)
		member.PUT("/projects/:id", projects.UpdateProject)
		member.DELETE("/projects/:id", projects.DeleteProject)
		member.POST("/upload", projects.UploadImage)
		member.POST("/education", education.CreateEducation)
		member.DELETE("/education/:id", education.DeleteEducation)
		member.POST("/experience", experience.CreateExperience)
		member.DELETE("/experience/:id", experience.DeleteExperience)

	}

	return r
}
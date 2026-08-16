package routes

import (
	"os"
	"strings"

	"github.com/Piyadanai03/portfolio-api/controllers/achievement"
	"github.com/Piyadanai03/portfolio-api/controllers/auth"
	"github.com/Piyadanai03/portfolio-api/controllers/education"
	"github.com/Piyadanai03/portfolio-api/controllers/experience"
	"github.com/Piyadanai03/portfolio-api/controllers/portfolio"
	"github.com/Piyadanai03/portfolio-api/controllers/profile"
	"github.com/Piyadanai03/portfolio-api/controllers/projects"
	"github.com/Piyadanai03/portfolio-api/controllers/technologies"
	_ "github.com/Piyadanai03/portfolio-api/docs"
	"github.com/Piyadanai03/portfolio-api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := []string{}

	if originsEnv != "" {
		allowedOrigins = strings.Split(originsEnv, ",")
	} else {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.HEAD("/ping", func(c *gin.Context) {
		c.Status(200)
	})
	v1 := r.Group("/api/v1")

	v1.GET("/projects", projects.GetProjects)
	v1.GET("/projects/:id", projects.GetProjectByID)
	v1.GET("/experiences", experience.GetExperiences)
	v1.GET("/home", portfolio.GetHomeData)
	v1.GET("/about", portfolio.GetAbout)
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

		member.GET("/education", education.GetEducations)
		member.GET("/education/:id", education.GetEducationByID)
		member.POST("/education", education.CreateEducation)
		member.PUT("/education/:id", education.UpdateEducation)
		member.DELETE("/education/:id", education.DeleteEducation)

		member.GET("/experiences", experience.GetExperiences)
		member.GET("/experience/:id", experience.GetExperienceByID)
		member.POST("/experience", experience.CreateExperience)
		member.PUT("/experience/:id", experience.UpdateExperience)
		member.DELETE("/experience/:id", experience.DeleteExperience)

		member.GET("/tech", technologies.GetTechnologies)
		member.GET("/tech/:id", technologies.GetTechnologyByID)
		member.POST("/tech", technologies.CreateTech)
		member.PUT("/tech/:id", technologies.UpdateTech)
		member.DELETE("/tech/:id", technologies.DeleteTech)

		member.GET("/profile", profile.GetProfile)
		member.PUT("/profile", profile.UpdateProfile)

		member.GET("/achievement", achievement.GetAchievement)
		member.GET("/achievement/:id", achievement.GetAchievementByID)
		member.POST("/achievement", achievement.CreateAchievement)
		member.PUT("/achievement/:id", achievement.UpdateAchievement)
		member.DELETE("/achievement/:id", achievement.DeleteAchievement)
	}

	return r
}

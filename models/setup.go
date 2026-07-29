package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username        string    `gorm:"unique;not null" json:"username"`
	PasswordHash    string    `gorm:"not null" json:"-"`
	FullName        string    `json:"fullName"`
	Position        string    `json:"position"`
	BioText         string    `json:"bioText"`
	Address         string    `json:"address"`
	ProfileImageURL string    `json:"profileImageURL"`
	ResumeURL       string    `json:"resumeURL"`
	CreatedAt       time.Time `json:"createdAt"`

	Projects    []Project    `gorm:"foreignKey:UserID" json:"projects,omitempty"`
	Experiences []Experience `gorm:"foreignKey:UserID" json:"experiences,omitempty"`
	Studies     []Study      `gorm:"foreignKey:UserID" json:"studies,omitempty"`
	Contacts    []Contact    `gorm:"foreignKey:UserID" json:"contacts"`
}

type Project struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	ExperienceID  *uuid.UUID `gorm:"type:uuid" json:"experienceID"`
	Title         string     `gorm:"not null" json:"title"`
	Description   string     `json:"description"`
	CoverImageURL string     `json:"coverImageURL"`
	GithubURL     string     `json:"githubURL"`
	CreatedAt     time.Time  `json:"createdAt"`

	Images       []ProjectImage `gorm:"foreignKey:ProjectID" json:"images"`
	Technologies []Technology   `gorm:"many2many:project_technologies;" json:"technologies"`
	Experience   *Experience    `gorm:"foreignKey:ExperienceID" json:"experience"`
}

type ProjectImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null" json:"projectID"`
	ImageURL  string    `gorm:"not null" json:"imageURL"`
	Caption   string    `json:"caption"`
}

type Technology struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `gorm:"unique;not null" json:"name"`
	Category string    `json:"category"`
	IconURL  string    `json:"iconURL"`
}

type Experience struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	JobTitle    string     `json:"jobTitle"`
	Company     string     `json:"company"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"` // ใช้ pointer เพื่อให้เป็น NULL ได้กรณีปัจจุบันยังทำอยู่
	Description string     `json:"description"`

	Projects []Project `gorm:"foreignKey:ExperienceID" json:"projects,omitempty"`
}

type Study struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Degree         string    `json:"degree"`
	Faculty        string    `json:"faculty"`
	Major          string    `json:"major"`
	Institution    string    `json:"institution"`
	GPA            *float64  `json:"gpa"`
	GraduationDate time.Time `json:"graduationDate"`
}

type Contact struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id,omitempty"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	PlatformName string    `json:"platformName"`
	URLValue     string    `json:"urlValue"`
	IconURL      string    `json:"iconURL"`
	IsActive     *bool     `gorm:"default:true" json:"isActive"`
}

type Achievement struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	
	Title        string     `json:"title"`
	Category     string     `json:"category"`
	DateAchieved time.Time  `json:"dateAchieved"`
	ProjectID    *uuid.UUID `json:"projectID"`
}

// ฟังก์ชันสั่ง Run Migration
func MigrateDB(db *gorm.DB) {
	
	db.AutoMigrate(
		&User{},
		&Project{},
		&ProjectImage{},
		&Technology{},
		&Experience{},
		&Study{},
		&Contact{},
		&Achievement{},
	)
}


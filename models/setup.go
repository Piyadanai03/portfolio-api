package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// 1. ตาราง users (Admin & Profile)
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
	
	// Relationships
	Projects    []Project    `gorm:"foreignKey:UserID" json:"projects,omitempty"`
	Experiences []Experience `gorm:"foreignKey:UserID" json:"experiences,omitempty"`
	Studies     []Study      `gorm:"foreignKey:UserID" json:"studies,omitempty"`
	Contacts    []Contact    `gorm:"foreignKey:UserID" json:"contacts"`
}

type Project struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Title         string    `gorm:"not null" json:"title"`
	Description   string    `json:"description"`
	CoverImageURL string    `json:"coverImageURL"`
	GithubURL     string    `json:"githubURL"`
	CreatedAt     time.Time `json:"createdAt"`
	Images       []ProjectImage `gorm:"foreignKey:ProjectID" json:"images"`
	Technologies []Technology   `gorm:"many2many:project_technologies;" json:"technologies"`
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

// 5. ตาราง experiences
type Experience struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	JobTitle    string
	Company     string
	StartDate   time.Time
	EndDate     *time.Time // ใช้ pointer เพื่อให้เป็น NULL ได้กรณีปัจจุบันยังทำอยู่
	Description string
}

// 6. ตาราง study
type Study struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	Degree         string
	Major          string
	Institution    string
	GraduationDate time.Time
}

// 7. ตาราง contact_info
type Contact struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id,omitempty"`
	UserID       uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	PlatformName string    `json:"platformName"`
	URLValue     string    `json:"urlValue"`
	IconURL      string    `json:"iconURL"`
	IsActive     *bool      `gorm:"default:true" json:"isActive"`
}

// 8. ตาราง achievements
type Achievement struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null"`
	ProjectID    *uuid.UUID `gorm:"type:uuid"` // เชื่อมโปรเจกต์ (ถ้ามี)
	Title        string
	Category     string // award หรือ training
	DateAchieved time.Time
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

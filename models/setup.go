package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 1. ตาราง users (Admin & Profile)
type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username        string    `gorm:"unique;not null"`
	PasswordHash    string    `gorm:"not null"`
	FullName        string
	BioText         string
	Address         string
	ProfileImageURL string
	ResumeURL       string
	CreatedAt       time.Time
	// Relationships
	Projects    []Project    `gorm:"foreignKey:UserID"`
	Experiences []Experience `gorm:"foreignKey:UserID"`
	Studies     []Study      `gorm:"foreignKey:UserID"`
	Contacts    []Contact    `gorm:"foreignKey:UserID"`
}

// 2. ตาราง projects
type Project struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Title         string    `gorm:"not null"`
	Description   string
	CoverImageURL string
	GithubURL     string
	CreatedAt     time.Time
	// Relationships
	Images       []ProjectImage `gorm:"foreignKey:ProjectID"`
	Technologies []Technology   `gorm:"many2many:project_technologies;"`
}

// 3. ตาราง project_images
type ProjectImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null"`
	ImageURL  string    `gorm:"not null"`
	Caption   string
}

// 4. ตาราง technologies
type Technology struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `gorm:"unique;not null"`
	Category string    // Backend, Frontend, AI, etc.
	IconURL  string
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
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	PlatformName string
	URLValue     string
	IconURL      string
	IsActive     bool `gorm:"default:true"`
}

// 8. ตาราง achievements
type Achievement struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
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
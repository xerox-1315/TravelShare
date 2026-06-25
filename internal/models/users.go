package models

import "time"

type User struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Username       string    `json:"username" gorm:"unique;not null;size:50"`
	Email          string    `json:"email" gorm:"unique;not null;size:100"`
	Password       string    `json:"-" gorm:"not null;size:100"`
	ProfileImage   string    `json:"profile_image"`
	Description    string    `json:"description"`
	PointsCount    int       `json:"points_count" gorm:"default:0"`
	LifehacksCount int       `json:"lifehacks_count" gorm:"default:0"`
	TravelLevel    int       `json:"travel_level" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

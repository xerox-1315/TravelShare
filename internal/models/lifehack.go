package models

import "time"

type Lifehack struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	RegionID    uint      `json:"region_id"`
	SeasonID    uint      `json:"season_id"`
	CategoryID  uint      `json:"category_id"`
	Title       string    `json:"title" gorm:"not null;size:200"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommentLifehack struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	LifehackID uint   `json:"lifehack_id" gorm:"not null"`
	Text       string `json:"text"`
	Image      string `json:"image"`
}

type VoteLifehack struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	LifehackID uint   `json:"lifehack_id" gorm:"not null"`
	TypeVote   string `json:"type" gorm:"not null;size:50"`
}

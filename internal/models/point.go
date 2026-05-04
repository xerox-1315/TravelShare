package models

import "time"

type Point struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	RegionID       uint      `json:"region_id"`
	SeasonID       uint      `json:"season_id"`
	CategoryID     uint      `json:"category_id"`
	DifficultLevel int       `json:"difficult_level"`
	Title          string    `json:"title" gorm:"not null;size:200"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CommentPoint struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  uint   `json:"user_id" gorm:"not null"`
	PointID uint   `json:"point_id" gorm:"not null"`
	Text    string `json:"text"`
	Image   string `json:"image"`
}

type VotePoint struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id" gorm:"not null"`
	PointID  uint   `json:"point_id" gorm:"not null"`
	TypeVote string `json:"type" gorm:"not null;size:50"`
}

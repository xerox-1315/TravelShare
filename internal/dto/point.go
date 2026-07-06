package dto

import "time"

type PointResponse struct {
	ID        uint    `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	SeasonID  *uint   `json:"season_id"`
}

type PointFullResponse struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	RegionID       *uint     `json:"region_id"`
	SeasonID       *uint     `json:"season_id"`
	CategoryID     *uint     `json:"category_id"`
	DifficultLevel int       `json:"difficult_level"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreatePointRequest struct {
	Latitude       float64 `json:"latitude" binding:"required"`
	Longitude      float64 `json:"longitude" binding:"required"`
	RegionID       *uint   `json:"region_id"`
	SeasonID       *uint   `json:"season_id"`
	CategoryID     *uint   `json:"category_id"`
	DifficultLevel int     `json:"difficult_level"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description" binding:"required"`
	Image          string  `json:"image"`
}

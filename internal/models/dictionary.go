package models

type Region struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null;size:200"`
	Description string `json:"description"`
}

type Season struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"not null;size:50"`
}

type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"not null;size:100"`
	Description string `json:"description"`
}

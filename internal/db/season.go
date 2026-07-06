package db

import (
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/models"
	"gorm.io/gorm"
)

type SeasonRepository struct {
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (sr *SeasonRepository) GetSeasonByID(id uint) (string, error) {
	var season models.Season
	result := sr.db.First(&season, id)
	if result.Error != nil {
		return "", errs.ErrorSeasonNotFound
	}
	return season.Title, nil
}

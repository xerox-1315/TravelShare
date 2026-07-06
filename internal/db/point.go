package db

import (
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/models"
	"gorm.io/gorm"
)

type PointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) *PointRepository {
	return &PointRepository{db: db}
}

func (pr *PointRepository) GetAllPoints() ([]models.Point, error) {
	var points []models.Point
	result := pr.db.Raw(`
        SELECT 
            id, user_id, region_id, season_id, category_id,
            difficult_level, title, description, image,
            ST_X(coords::geometry) as longitude,
            ST_Y(coords::geometry) as latitude
        FROM points
    `).Scan(&points)
	if result.Error != nil {
		return nil, errs.ErrorFindAllPoints
	}
	return points, nil
}

func (pr *PointRepository) GetPointByID(id int) (*models.Point, error) {
	var point models.Point
	result := pr.db.Raw(`
	SELECT id, user_id, region_id, season_id, category_id,
            difficult_level, title, description, image,
            ST_X(coords::geometry) as longitude,
            ST_Y(coords::geometry) as latitude
			FROM points WHERE id = ?`, id).Scan(&point)
	if result.Error != nil || point.ID == 0 {
		return nil, errs.ErrorFindPointByID
	}
	return &point, nil
}

func (pr *PointRepository) CreatePoint(point *models.Point) error {
	// обычным SQL запросом создаем новую метку в БД
	result := pr.db.Exec(`
		INSERT INTO points (coords, user_id, region_id, season_id,
		category_id, difficult_level, title, description, image) VALUES
		(ST_SetSRID(ST_MakePoint(?, ?), 4326), ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`,
		point.Longitude, point.Latitude, point.UserID, point.RegionID,
		point.SeasonID, point.CategoryID, point.DifficultLevel, point.Title,
		point.Description, point.Image)

	if result.Error != nil {
		return errs.ErrorCreatePoint
	}
	return nil
}

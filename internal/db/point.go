package db

import (
	"errors"

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
	if result.Error != nil {
		return nil, errs.ErrorFindPointByID
	}
	if point.ID == 0 {
		return nil, errs.ErrorNotFoundPointByID
	}
	return &point, nil
}

func (pr *PointRepository) GetPointsNearby(latitude, longitude float64, radius int) ([]*models.Point, error) {
	var points []*models.Point
	result := pr.db.Raw(`
	SELECT id, user_id, season_id, region_id, category_id,
           difficult_level, title, description, image,
           ST_X(coords::geometry) as longitude,
           ST_Y(coords::geometry) as latitude
    FROM points
    WHERE ST_DWithin(
        coords,
        ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
        ?
    )`,
		longitude, latitude, radius).Scan(&points)
	if result.Error != nil {
		return nil, errs.ErrorFindPointsNearby
	}
	return points, nil
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

func (pr *PointRepository) GetUserVoteByPointID(userID, pointID uint) (bool, *models.VotePoint, error) {
	var votePoint models.VotePoint
	result := pr.db.Where("user_id = ?", userID).Where("point_id = ?", pointID).First(&votePoint)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil, nil // не голосовал
		}
		return false, nil, result.Error
	}
	return true, &votePoint, nil // уже голосовал

}

func (pr *PointRepository) SetVote(userID, pointID uint, typeVote string) (*models.VotePoint, error) {
	votePoint := models.VotePoint{
		UserID:   userID,
		PointID:  pointID,
		TypeVote: typeVote,
	}
	result := pr.db.Create(&votePoint)
	if result.Error != nil {
		return nil, errs.ErrorCreateVotePoint
	}
	return &votePoint, nil
}

func (pr *PointRepository) UpdateVote(userID, pointID uint, typeVote string) (*models.VotePoint, error) {
	votePoint := models.VotePoint{
		UserID:  userID,
		PointID: pointID,
	}
	// обновить оценку на метку
	result := pr.db.Model(&votePoint).Where("user_id = ? AND point_id = ?", userID, pointID).Update("type_vote", typeVote)
	if result.Error != nil {
		return nil, result.Error
	}
	return &votePoint, nil
}

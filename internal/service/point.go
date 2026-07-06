package service

import (
	"github.com/xerox-1315/TravelShare.git/internal/db"
	"github.com/xerox-1315/TravelShare.git/internal/dto"
	"github.com/xerox-1315/TravelShare.git/internal/models"
)

// логический слой работы с метками
type PointService struct {
	repo *db.PointRepository
}

// инициализация слоя
func NewPointService(repo *db.PointRepository) *PointService {
	return &PointService{repo: repo}
}

// получение всех доступных меток
func (ps *PointService) GetAllPoints() ([]dto.PointResponse, error) {
	// достаем метки из БД (слой работы с данными)
	points, err := ps.repo.GetAllPoints()
	if err != nil {
		return nil, err
	}

	// массив всех DTO точек для передачи в хендлер
	var pointsDtoList []dto.PointResponse
	// итерируемся по всем точкам
	for _, point := range points {
		// создаем DTO
		pointDTO := dto.PointResponse{
			ID:        point.ID,
			Latitude:  point.Latitude,
			Longitude: point.Longitude,
			SeasonID:  point.SeasonID,
		}
		// добавляем в массив
		pointsDtoList = append(pointsDtoList, pointDTO)
	}
	return pointsDtoList, nil
}

func (ps *PointService) GetPointByID(id int) (*dto.PointFullResponse, error) {
	point, err := ps.repo.GetPointByID(id)
	if err != nil {
		return nil, err
	}
	pointDTO := dto.PointFullResponse{
		ID:             point.ID,
		UserID:         point.UserID,
		Latitude:       point.Latitude,
		Longitude:      point.Longitude,
		RegionID:       point.RegionID,
		SeasonID:       point.SeasonID,
		CategoryID:     point.CategoryID,
		DifficultLevel: point.DifficultLevel,
		Title:          point.Title,
		Description:    point.Description,
		Image:          point.Image,
		CreatedAt:      point.CreatedAt,
		UpdatedAt:      point.UpdatedAt,
	}
	return &pointDTO, nil
}

func (ps *PointService) CreatePoint(userID uint, req dto.CreatePointRequest) (*dto.PointFullResponse, error) {
	// создаем структуру точки для создания в БД
	point := &models.Point{
		UserID:         userID,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		RegionID:       req.RegionID,
		SeasonID:       req.SeasonID,
		CategoryID:     req.CategoryID,
		DifficultLevel: req.DifficultLevel,
		Title:          req.Title,
		Description:    req.Description,
		Image:          req.Image,
	}
	err := ps.repo.CreatePoint(point)
	if err != nil {
		return nil, err
	}
	response := dto.PointFullResponse{
		ID:             point.ID,
		UserID:         userID,
		Latitude:       point.Latitude,
		Longitude:      point.Longitude,
		RegionID:       point.RegionID,
		SeasonID:       point.SeasonID,
		CategoryID:     point.CategoryID,
		DifficultLevel: point.DifficultLevel,
		Title:          point.Title,
		Description:    point.Description,
		Image:          point.Image,
	}
	return &response, nil
}

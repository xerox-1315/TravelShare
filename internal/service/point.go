package service

import (
	"errors"

	"github.com/xerox-1315/TravelShare.git/errs"
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

func (ps *PointService) GetPointsNearby(latitude, longitude float64, radius int) ([]*dto.PointResponse, error) {
	points, err := ps.repo.GetPointsNearby(latitude, longitude, radius)
	if err != nil {
		return nil, err
	}
	var pointsDTO []*dto.PointResponse
	for _, point := range points {
		pointDTO := dto.PointResponse{
			ID:        point.ID,
			Latitude:  point.Latitude,
			Longitude: point.Longitude,
			SeasonID:  point.SeasonID,
		}
		pointsDTO = append(pointsDTO, &pointDTO)
	}
	if len(pointsDTO) == 0 {
		return nil, errs.ErrorNotFoundPointsNearby
	}
	return pointsDTO, nil
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

func (ps *PointService) SetVote(userID, pointID uint, typeVote string) (*models.VotePoint, error) {
	// проверка голоса на один из двух типов
	if typeVote != "like" && typeVote != "dislike" {
		return nil, errs.ErrorInvalidTypeOfVote
	}

	// проверка существования метки с данным ID
	_, err := ps.repo.GetPointByID(int(pointID))
	if errors.Is(err, errs.ErrorNotFoundPointByID) {
		return nil, err
	}

	// проверка, ставил ли данный пользователь голос на данную метку
	exist, voteExist, err := ps.repo.GetUserVoteByPointID(userID, pointID)
	if err != nil {
		return nil, err
	}
	if exist {
		// если новая метка отличная от существующей
		if typeVote != voteExist.TypeVote {
			// обновляем данные в БД на новые
			voteDTO, err := ps.repo.UpdateVote(userID, pointID, typeVote)
			if err != nil {
				return nil, err
			}
			return voteDTO, nil
		} else {
			// иначе не даем поставить еще такую же оценку
			return nil, errs.ErrorExistVoteFromUser
		}
	}

	// случай, когда пользователь еще не ставил оценку на данную точку
	voteDTO, err := ps.repo.SetVote(userID, pointID, typeVote)
	if err != nil {
		return nil, err
	}
	return voteDTO, nil
}

func (ps *PointService) UpdateVote(userID, pointID uint, typeVote string) (*models.VotePoint, error) {
	// проверка голоса на один из двух типов
	if typeVote != "like" && typeVote != "dislike" {
		return nil, errs.ErrorInvalidTypeOfVote
	}

	// проверка существования метки с данным ID
	_, err := ps.repo.GetPointByID(int(pointID))
	if errors.Is(err, errs.ErrorNotFoundPointByID) {
		return nil, err
	}

	// проверка, ставил ли данный пользователь голос на данную метку
	exist, _, err := ps.repo.GetUserVoteByPointID(userID, pointID)
	if err != nil {
		return nil, err
	}
	if !exist {
		// если голоса от пользователя на данную метку не было
		return nil, errs.ErrorNotExistVoteFromUser
	}

	// обновляем данные в БД на новые
	voteDTO, err := ps.repo.UpdateVote(userID, pointID, typeVote)
	if err != nil {
		return nil, err
	}
	return voteDTO, nil
}

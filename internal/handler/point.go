package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/dto"
	"github.com/xerox-1315/TravelShare.git/internal/service"
)

// слой обрабаотки запросов с метками
type PointHandler struct {
	service *service.PointService
}

// инициализация слоя
func NewPointHandler(service *service.PointService) *PointHandler {
	return &PointHandler{service: service}
}

func (ph *PointHandler) GetAllPoints(c *gin.Context) {
	// получение всех меток
	pointsDTO, err := ph.service.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	c.JSON(http.StatusOK, pointsDTO)
}

func (ph *PointHandler) GetPointByID(c *gin.Context) {
	// достаем id из параметров запроса
	id := c.Param("id")
	// если параметр не передан
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не передан параметр запроса"})
		return
	}
	// конвертация id из строкового типа в числовой
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидный параметр запроса"})
		return
	}
	// получаем метку по ID
	pointDTO, err := ph.service.GetPointByID(idInt)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFoundPointByID) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Метка с таким ID не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	c.JSON(http.StatusOK, pointDTO)
}

func (ph *PointHandler) GetPointsNearby(c *gin.Context) {
	// получение параметров запроса (широта, долгота, радиус)
	latitudeStr := c.Query("lat")
	longitudeStr := c.Query("lng")
	radiusStr := c.Query("radius")
	if latitudeStr == "" || longitudeStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидные параметры запроса"})
		return
	}

	// переводим из строки в число с плавающей точкой
	latitude, err1 := strconv.ParseFloat(latitudeStr, 64)
	longitude, err2 := strconv.ParseFloat(longitudeStr, 64)
	// радиус переводим в целое число
	radius, err3 := strconv.Atoi(radiusStr)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидные параметры запроса"})
		return
	}

	// обращаемся к service и получаем все метки в области
	pointsDTO, err := ph.service.GetPointsNearby(latitude, longitude, radius)
	if err != nil {
		if errors.Is(err, errs.ErrorNotFoundPointsNearby) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Метки в данном диапазоне не найдены"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	c.JSON(http.StatusOK, pointsDTO)
}

func (ph *PointHandler) CreatePoint(c *gin.Context) {
	var req dto.CreatePointRequest
	// парсим тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы неккоректные данные"})
		return
	}

	// парсим user_id из middleware
	userID, _ := c.Get("user_id")
	response, err := ph.service.CreatePoint(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (ph *PointHandler) SetVote(c *gin.Context) {
	var req dto.VotePoint
	// парсим тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы неккоректные данные"})
		return
	}

	// парсим user_id из middleware
	userID, _ := c.Get("user_id")
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
	}
	// устанавливаем оценку в service
	votePoint, err := ph.service.SetVote(id, req.PointID, req.TypeVote)
	if err != nil {
		if errors.Is(err, errs.ErrorInvalidTypeOfVote) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Передан неверный тип оценки (ожидается like/dislike)"})
			return
		} else if errors.Is(err, errs.ErrorNotFoundPointByID) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Передан несуществующий идентификатор метки"})
			return
		} else if errors.Is(err, errs.ErrorExistVoteFromUser) {
			c.JSON(http.StatusConflict, gin.H{"message": "Такая оценка от данного пользователя уже существует"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
	}
	c.JSON(http.StatusCreated, &votePoint)
}

func (ph *PointHandler) UpdateVote(c *gin.Context) {
	var req dto.VotePoint
	// парсим тело запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы неккоректные данные"})
		return
	}

	// парсим user_id из middleware
	userID, _ := c.Get("user_id")
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
	}
	// обновляем оценку в service
	votePoint, err := ph.service.UpdateVote(id, req.PointID, req.TypeVote)
	if err != nil {
		if errors.Is(err, errs.ErrorInvalidTypeOfVote) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Передан неверный тип оценки (ожидается like/dislike)"})
			return
		} else if errors.Is(err, errs.ErrorNotFoundPointByID) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Передан несуществующий идентификатор метки"})
			return
		} else if errors.Is(err, errs.ErrorNotExistVoteFromUser) {
			c.JSON(http.StatusConflict, gin.H{"message": "Данный пользователь еще не оставлял оценку на эту метку"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
	}
	c.JSON(http.StatusOK, &votePoint)
}

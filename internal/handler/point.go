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
	pointsDTO, err := ph.service.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	c.JSON(http.StatusOK, pointsDTO)
}

func (ph *PointHandler) GetPointByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не передан параметр запроса"})
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидный параметр запроса"})
		return
	}
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
	latitudeStr := c.Query("lat")
	longitudeStr := c.Query("lng")
	radiusStr := c.Query("radius")
	if latitudeStr == "" || longitudeStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидные параметры запроса"})
		return
	}

	latitude, err1 := strconv.ParseFloat(latitudeStr, 64)
	longitude, err2 := strconv.ParseFloat(longitudeStr, 64)
	radius, err3 := strconv.Atoi(radiusStr)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Невалидные параметры запроса"})
		return
	}
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

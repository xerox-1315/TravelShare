package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/service"
)

// слой обработки запросов
type UserHandler struct {
	service *service.UserService
}

// инициализация слоя
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// структура запроса регистрации пользователя
type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required"`
	ProfileImage string `json:"profile_image"`
	Description  string `json:"description"`
}

// структура запроса входа пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// распаковываем JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		// при ошибке возвращаем, что данные переданы неправильно
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы некорректные данные запроса"})
		return
	}

	// выполняем регистрацию пользователя
	err := uh.service.Register(req.Username, req.Email, req.Password, req.ProfileImage, req.Description)
	// обработка ошибок
	if err != nil {
		switch err {
		case errs.ErrorUserWithEmailExist, errs.ErrorUserWithUsernameExist:
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		case errs.ErrorPasswordLenght:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		}
		return
	}
	// успешная регистрация
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы неккоректные данные запроса"})
		return
	}
	// выполняем вход пользователя и получаем JWT-токен
	token, err := uh.service.Login(req.Email, req.Password)
	// обработка ошибок
	if err != nil {
		switch err {
		case errs.ErrorUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case errs.ErrorWrongPassword:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		}
		return
	}
	// успешная аутентификация - передаем токен на выход
	c.JSON(http.StatusOK, gin.H{"token": token})
}

package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/dto"
	"github.com/xerox-1315/TravelShare.git/internal/service"
)

// слой обработки запросов (пользователи)
type UserHandler struct {
	service *service.UserService
}

// инициализация слоя
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (uh *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

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
	// код отправлен, можно ожидать верификации
	c.JSON(http.StatusOK, gin.H{"message": "Код подтверждения отправлен на почту"})
}

func (uh *UserHandler) Verify(c *gin.Context) {
	var req dto.VerifyRequest
	// распаковываем JSON из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		// при ошибке возвращаем, что данные переданы неправильно
		c.JSON(http.StatusBadRequest, gin.H{"message": "Переданы некорректные данные запроса"})
		return
	}

	user, err := uh.service.Verify(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, errs.ErrorInvalidCode) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный код подтверждения"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
	}
	c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

func (uh *UserHandler) GetProfile(c *gin.Context) {
	// получаем userID из middleware
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "не авторизован"})
		return
	}
	user, err := uh.service.GetProfile(userID.(uint))
	if err != nil {
		// пользователь с таким ID не найден
		c.JSON(http.StatusNotFound, gin.H{"message": "Такой пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, user)
}

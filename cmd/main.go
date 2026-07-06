package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xerox-1315/TravelShare.git/internal/db"
	"github.com/xerox-1315/TravelShare.git/internal/handler"
	"github.com/xerox-1315/TravelShare.git/internal/service"
)

func main() {
	// подгрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки файла .env")
		return
	}
	// создаем подключение к БД
	database, err := db.NewConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// инициализация слоев системы (пользователи)
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// инициализация слоев системы (метки)
	pointRepo := db.NewPointRepository(database)
	pointService := service.NewPointService(pointRepo)
	pointHandler := handler.NewPointHandler(pointService)

	// инициализируем роутер
	r := gin.Default()

	// роуты авторизации - публичные
	auth := r.Group("/auth")
	{
		// регистрация
		auth.POST("/register", userHandler.Register)
		// вход
		auth.POST("/login", userHandler.Login)
	}

	// защищенные роуты
	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/points", pointHandler.GetAllPoints)
		protected.POST("/point", pointHandler.CreatePoint)
		protected.GET("/point/:id", pointHandler.GetPointByID)
	}

	// запускаем сервер
	r.Run(":8080")
}

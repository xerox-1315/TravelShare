package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xerox-1315/TravelShare.git/internal/cache"
	"github.com/xerox-1315/TravelShare.git/internal/db"
	"github.com/xerox-1315/TravelShare.git/internal/grpc"
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

	// инициализируем gRPC клиент
	emailClient, err := grpc.NewEmailClient(os.Getenv("EMAIL_SERVICE_ADDR"))
	if err != nil {
		log.Fatal("Ошибка подключения к Email Service: ", err)
	}
	// инициализируем redis клиент
	redisClient := cache.NewRedisClient(os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))

	// инициализация слоев системы (пользователи)
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo, emailClient, redisClient)
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
		// верификация
		auth.POST("/verify", userHandler.Verify)
	}

	// защищенные роуты
	protected := r.Group("/")
	protected.Use(handler.AuthMiddleware())
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/points", pointHandler.GetAllPoints)
		protected.POST("/point", pointHandler.CreatePoint)
		protected.GET("/point/:id", pointHandler.GetPointByID)
		protected.GET("/points/nearby", pointHandler.GetPointsNearby)
		protected.POST("/vote", pointHandler.SetVote)
		protected.PATCH("/vote", pointHandler.UpdateVote)
	}

	// запускаем сервер
	r.Run(":8080")
}

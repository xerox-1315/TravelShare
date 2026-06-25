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
	// инициализация слоев системы
	userRepo := db.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// инициализируем роутер
	r := gin.Default()

	// роуты авторизации
	auth := r.Group("/auth")
	{
		// регистрация
		auth.POST("/register", userHandler.Register)
		// вход
		auth.POST("/login", userHandler.Login)
	}

	// запускаем сервер
	r.Run(":8080")
}

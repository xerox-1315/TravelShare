package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xerox-1315/TravelShare.git/internal/service"
)

// промежуточная оболочка для авторизации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// достаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Токен отсутствует"})
			c.Abort()
			return
		}

		// убираем "Bearer " и отсавляем только токен
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный формат токена"})
			c.Abort()
			return
		}

		// парсим токен
		claims, err := service.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный токен"})
			c.Abort()
			return
		}

		// устанавливаем user_id в заголовки для хендлеров
		c.Set("user_id", claims.UserID)
		// идем в хендлер
		c.Next()
	}
}

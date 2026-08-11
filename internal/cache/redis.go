package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port, password string) *RedisClient {
	// создание redis клиента
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return &RedisClient{client: client}
}

// сохранение данных пользователя
func (r *RedisClient) SaveUser(email string, data []byte) error {
	return r.client.Set(context.Background(), "keeper"+email, data, 10*time.Minute).Err()
}

// получение данных пользователя по email
func (r *RedisClient) GetUser(email string) ([]byte, error) {
	return r.client.Get(context.Background(), "keeper"+email).Bytes()
}

// удаление данных пользователя
func (r *RedisClient) DelUser(email string) error {
	return r.client.Del(context.Background(), "keeper"+email).Err()
}

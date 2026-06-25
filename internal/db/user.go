package db

import (
	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/models"
	"gorm.io/gorm"
)

// слой работы с данными
type UserRepository struct {
	// подключение ORM
	db *gorm.DB
}

// инициализация слоя
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// создание пользователя
func (ur *UserRepository) CreateUser(user *models.User) error {
	// проверка юзернейма на уникальность
	result := ur.db.Create(user)
	if result.Error != nil {
		return errs.ErrorCreateUser
	}
	return nil
}

// поиск пользователя по email
func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// пытаемся найти первого пользователя с данным email
	result := ur.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errs.ErrorUserNotFound
	}
	return &user, nil
}

// поиск пользоваетля по username
func (ur *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	// пытаемся найти первого пользователя с данным ником
	result := ur.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, errs.ErrorUserNotFound
	}
	return &user, nil
}

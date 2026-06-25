package service

import (
	"unicode/utf8"

	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/db"
	"github.com/xerox-1315/TravelShare.git/internal/models"
)

// слой логики приложения
type UserService struct {
	repo *db.UserRepository
}

// инициализация слоя
func NewUserService(repo *db.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Register(username, email, password, profileImage, description string) error {
	// проверка на уникальность email
	_, err := us.repo.GetUserByEmail(email)
	if err == nil {
		return errs.ErrorUserWithEmailExist
	}

	// проверка на уникальность никнейма
	_, err = us.repo.GetUserByUsername(username)
	if err == nil {
		return errs.ErrorUserWithUsernameExist
	}

	// проверка длины пароля
	if utf8.RuneCountInString(password) < 8 {
		return errs.ErrorPasswordLenght
	}

	// создание экземпляра пользователя с переданными данными
	user := &models.User{
		Username:     username,
		Email:        email,
		Password:     password,
		ProfileImage: profileImage,
		Description:  description,
	}
	return us.repo.CreateUser(user)
}

func (us *UserService) Login(email, password string) (string, error) {
	// находим пользователя в БД
	user, err := us.repo.GetUserByEmail(email)
	if err != nil {
		return "", errs.ErrorUserNotFound
	}
	// проверяем пароли на совпадение
	if password != user.Password {
		return "", errs.ErrorWrongPassword
	}
	// создаем JWT-токен для хранения сессии пользователя
	token, err := CreateToken(user.ID)
	if err != nil {
		return "", errs.ErrorGenerateToken
	}
	return token, nil
}

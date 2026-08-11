package service

import (
	"encoding/json"
	"log"
	"unicode/utf8"

	"github.com/xerox-1315/TravelShare.git/errs"
	"github.com/xerox-1315/TravelShare.git/internal/cache"
	"github.com/xerox-1315/TravelShare.git/internal/db"
	"github.com/xerox-1315/TravelShare.git/internal/grpc"
	"github.com/xerox-1315/TravelShare.git/internal/models"
)

// слой логики приложения
type UserService struct {
	repo        *db.UserRepository
	emailClient *grpc.EmailClient
	redisClient *cache.RedisClient
}

// инициализация слоя
func NewUserService(repo *db.UserRepository, emailClient *grpc.EmailClient,
	redisClient *cache.RedisClient) *UserService {
	return &UserService{repo: repo, emailClient: emailClient, redisClient: redisClient}
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

	// упаковка данных пользователя в json-формат
	userData, _ := json.Marshal(map[string]string{
		"username":      username,
		"email":         email,
		"password":      password,
		"profile_image": profileImage,
		"description":   description,
	})
	// сохранение данных пользователя в redis хранилище
	err = us.redisClient.SaveUser(email, userData)
	if err != nil {
		log.Println("Ошибка сохранения в Redis:", err)
		return err
	}
	log.Println("Данные сохранены в Redis")

	// отправка кода на почту через email-сервис
	_, err = us.emailClient.SendVerificationCode(email)
	if err != nil {
		log.Println("Ошибка gRPC вызова:", err)
		return err
	}
	log.Println("Код отправлен через gRPC")
	return err
}

func (us *UserService) Verify(email, code string) (*models.User, error) {
	// проверка кода через email сервис
	verify, err := us.emailClient.VerifyCode(email, code)
	if err != nil {
		log.Println("ошибка проверка кода:", err)
		return nil, err
	}
	log.Println("код проверен через gRPC")
	if !verify {
		return nil, errs.ErrorInvalidCode
	}

	// если код верен - достаем данные пользователя из redis
	data, err := us.redisClient.GetUser(email)
	if err != nil {
		log.Println("ошибка получения данных из redis")
		return nil, err
	}
	log.Println("данные получены из redis")
	// распаковываем данные
	var userData map[string]string
	json.Unmarshal(data, &userData)
	// создание экземпляра пользователя с переданными данными
	user := &models.User{
		Username:     userData["username"],
		Email:        userData["email"],
		Password:     userData["password"],
		ProfileImage: userData["profileImage"],
		Description:  userData["description"],
	}
	return user, us.repo.CreateUser(user)
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

func (us *UserService) GetProfile(userID uint) (*models.User, error) {
	// обращаемся к БД для получения пользователя
	user, err := us.repo.GetProfile(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

package errs

import "errors"

var ErrorDBConnection = errors.New("Ошибка подключения к БД")
var ErrorCreateUser = errors.New("Ошибка создания пользователя")
var ErrorUserNotFound = errors.New("Такой пользователь не найден")
var ErrorUserWithEmailExist = errors.New("Пользователь с таким email уже существует")
var ErrorUserWithUsernameExist = errors.New("Пользователь с таким никнеймом уже существует")
var ErrorPasswordLenght = errors.New("Длина пароля должна быть не меньше 8 символов")
var ErrorGenerateToken = errors.New("Не удалось сгенерировать токен")
var ErrorInvalidToken = errors.New("Переданный токен невалиден")
var ErrorWrongPassword = errors.New("Неверный пароль")
var ErrorSeasonNotFound = errors.New("Сезон с таким id не найден")
var ErrorFindAllPoints = errors.New("Ошибка нахождения всех точек в БД")
var ErrorFindPointByID = errors.New("Ошибка нахождения точки в БД")
var ErrorCreatePoint = errors.New("Ошибка создания точки")
var ErrorFindPointsNearby = errors.New("Ошибка нахождения точек в области")
var ErrorNotFoundPointsNearby = errors.New("Метки в данном диапазоне не найдены")
var ErrorNotFoundPointByID = errors.New("Метка с таким ID не найдена")

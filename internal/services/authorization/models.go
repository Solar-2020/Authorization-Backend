package authorization

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/pkg/errors"
)

var (
	ErrorUserNotFound      = errors.New("Пользователь с указанным email не найден!")
	ErrorIncorrectPassword = errors.New("Указан неверный пароль!")
	ErrorNotValidCookie    = errors.New("Ваша кука больше недействительна, авторизуйтесь повторно")
)

type authorizationStorage interface {
	InsertPassword(pass models.Password) (err error)
	UpdatePassword(pass models.Password) (err error)
	SelectPasswordByUserID(userID int) (pass models.Password, err error)

	InsertCookie(cookie models.Cookie) (err error)
	SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error)
}

type accountClient interface {
	GetUserByUid(userID int) (user account.User, err error)
	GetUserByEmail(email string) (user account.User, err error)
	GetYandexUser(userToken string) (user account.User, err error)
	CreateUser(request account.User) (userID int, err error)
}

type errorWorker interface {
	NewError(httpCode int, responseError error, fullError error) (err error)
}

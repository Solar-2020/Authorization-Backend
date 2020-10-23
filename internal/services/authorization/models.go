package authorization

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
)

type accountClient interface {
	CreateUser(user models.Registration) (userID int, err error)
	GetUserIDByEmail(email string) (userID int, err error)
}

type authorizationStorage interface {
	InsertPassword(pass models.Password) (err error)
	UpdatePassword(pass models.Password) (err error)
	SelectPasswordByUserID(userID int) (pass models.Password, err error)

	InsertCookie(cookie models.Cookie) (err error)
	SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error)
}

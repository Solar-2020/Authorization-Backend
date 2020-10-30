package models

import (
	"github.com/Solar-2020/Account-Backend/pkg/models"
)

type AccountServiceInterface interface {
	GetUserByUid(userID int) (user models.User, err error)
	GetUserByEmail(email string) (user models.User, err error)
	CreateUser(request models.User) (userID int, err error)
}

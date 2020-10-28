package models

import (
	"github.com/Solar-2020/Account-Backend/pkg/api"
	"github.com/Solar-2020/Account-Backend/pkg/models"
)

type AccountServiceInterface interface {
	GetUserByUid(userID int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(request api.CreateRequest) (api.CreateResponse, error)
}

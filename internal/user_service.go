package internal

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"math/rand"
)

type AccountService struct {

}

func (s *AccountService) CreateUser(user models.Registration) (userID int, err error) {
	userID = rand.Intn(3000)
	return
}
func (s *AccountService) GetUserIDByEmail(email string) (userID int, err error) {
	//userID = rand.Intn(3000)
	userID = 2081
	return
}
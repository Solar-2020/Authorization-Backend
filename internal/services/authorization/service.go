package authorization

import (
	"database/sql"
	"errors"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"time"
)

type Service interface {
	Authorization(request models.Authorization) (cookie models.Cookie, err error)
	Registration(request models.Registration) (cookie models.Cookie, err error)
	GetUserIdByCookie(cookieValue string) (userID int, err error)
}

type service struct {
	authorizationStorage  authorizationStorage
	accountClient         accountClient
	defaultCookieLifeTime time.Duration
}

func NewService(authorizationStorage authorizationStorage, accountClient accountClient, defaultCookieLifeTime time.Duration) Service {
	return &service{
		authorizationStorage:  authorizationStorage,
		accountClient:         accountClient,
		defaultCookieLifeTime: defaultCookieLifeTime,
	}
}

func (s *service) Authorization(request models.Authorization) (cookie models.Cookie, err error) {
	userID, err := s.accountClient.GetUserIDByEmail(request.Login)
	if err != nil {
		return
	}

	err = s.checkLogoPass(userID, request.Password)
	if err != nil {
		return
	}

	cookie, err = s.createCookie(userID, s.defaultCookieLifeTime)
	if err != nil {
		return
	}

	err = s.authorizationStorage.InsertCookie(cookie)
	return

}

func (s *service) Registration(request models.Registration) (cookie models.Cookie, err error) {
	userID, err := s.accountClient.CreateUser(request)
	if err != nil {
		return
	}

	pass := s.generatePassword(userID, request.Password)

	err = s.authorizationStorage.InsertPassword(pass)
	if err != nil {
		return
	}

	cookie, err = s.createCookie(userID, s.defaultCookieLifeTime)
	if err != nil {
		return
	}

	err = s.authorizationStorage.InsertCookie(cookie)
	return
}

func (s *service) GetUserIdByCookie(cookieValue string) (userID int, err error) {
	cookie, err := s.authorizationStorage.SelectCookieByValue(cookieValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return userID, errors.New("Кука не действительна")
		}
		return
	}

	if cookie.Expiration.Before(time.Now()) {
		return userID, errors.New("Кука протухла")
	}

	return cookie.UserID, err
}

func (s *service) createCookie(userID int, expiration time.Duration) (cookie models.Cookie, err error) {

	return
}

func (s *service) checkLogoPass(userID int, userPassword string) (err error) {
	pass, err := s.authorizationStorage.SelectPasswordByUserID(userID)
	if err != nil {
		return
	}

	hashPassword := s.createPassword(userPassword, pass.Salt)
	if hashPassword != pass.HashPassword {
		return errors.New("Не верный пароль!")
	}

	return
}

func (s *service) generatePassword(userID int, userPassword string) (pass models.Password) {

	return
}

func (s *service) createPassword(userPassword string, salt string) (hashPass string) {

	return
}

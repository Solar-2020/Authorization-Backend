package authorization

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/Authorization-Backend/cmd/config"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/Solar-2020/GoUtils/common"
	"golang.org/x/crypto/pbkdf2"
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
}

func NewService(authorizationStorage authorizationStorage, accountClient accountClient) Service {
	return &service{
		authorizationStorage:  authorizationStorage,
		accountClient:         accountClient,
	}
}

func (s *service) Authorization(request models.Authorization) (cookie models.Cookie, err error) {
	//userID, err := s.accountClient.GetUserIDByEmail(request.Login)
	//if err != nil {
	//	return
	//}
	userID := request.Uid

	err = s.checkLogoPass(userID, request.Password)
	if err != nil {
		return
	}

	cookie, err = s.createCookie(userID, time.Duration(config.Config.DefaultCookieLifetime)*time.Second)
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

	cookie, err = s.createCookie(userID, time.Duration(config.Config.DefaultCookieLifetime)*time.Second)
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
	cookie.UserID = userID
	cookie.Expiration = time.Now().Add(expiration)
	cookie.Value = common.GetRandomString(int64(config.Config.SessionCookieLength),
		common.Capitals, common.Lowers, common.Numbers, common.LightSymbols)
	return
}

func (s *service) checkLogoPass(userID int, userPassword string) (err error) {
	basePass, err := s.authorizationStorage.SelectPasswordByUserID(userID)
	if err != nil {
		err = fmt.Errorf("not exist")
		return
	}

	hashPassword := s.hashPassword([]byte(userPassword), basePass.Salt)
	if !bytes.Equal(hashPassword, basePass.HashPassword) {
		err = errors.New("wrong password")
	}
	return
}

func (s *service) generatePassword(userID int, userPassword string) (pass models.Password) {
	pass.UserID = userID
	pass.Salt = md5.New().Sum([]byte(string(userID) + time.Now().String()))
	pass.HashPassword = s.hashPassword([]byte(userPassword), pass.Salt)
	pass.UpdateAt = time.Now()
	return
}

func (s *service) hashPassword(plainPassword []byte, salt []byte) (hashPass []byte) {
	return pbkdf2.Key(plainPassword, salt, 4096, 32, sha1.New)
}

package authorization

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"errors"
	models2 "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/Solar-2020/Authorization-Backend/cmd/config"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/Solar-2020/GoUtils/common"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/pbkdf2"
	"time"
)

type service struct {
	authorizationStorage authorizationStorage
	accountClient        accountClient
	errorWorker          errorWorker
}

func NewService(authorizationStorage authorizationStorage, accountClient accountClient, errorWorker errorWorker) *service {
	return &service{
		authorizationStorage: authorizationStorage,
		accountClient:        accountClient,
		errorWorker:          errorWorker,
	}
}

func (s *service) Authorization(request models.Authorization) (cookie models.Cookie, err error) {
	user, err := s.accountClient.GetUserByEmail(request.Login)
	if err != nil {
		return cookie, err
	}

	err = s.checkLogoPass(user.ID, request.Password)
	if err != nil {
		return
	}

	return s.addCookie(user.ID)
}

func (s *service) Registration(request models.Registration) (cookie models.Cookie, err error) {
	userID, err := s.createUser(request)
	if err != nil {
		return
	}

	pass := s.generatePassword(userID, request.Password)
	err = s.authorizationStorage.InsertPassword(pass)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	return s.addCookie(userID)
}

func (s *service) Yandex(userToken string) (cookie models.Cookie, err error) {
	user, err := s.accountClient.GetYandexUser(userToken)
	if err != nil {
		return
	}

	return s.addCookie(user.ID)
}

func (s *service) DublicateCookie(cookieValue string, lifetime int) (newCookie models.Cookie, err error) {
	userId, err := s.GetUserIdByCookie(cookieValue)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, errors.New("неверный исходный токен"), err)
		return
	}
	newCookie, err = s.addCookieWithLifetime(userId, time.Duration(lifetime))
	return
}

func (s *service) addCookie(userId int) (cookie models.Cookie, err error) {
	return s.addCookieWithLifetime(userId, time.Duration(config.Config.DefaultCookieLifetime))
}

func (s *service) addCookieWithLifetime(userId int, lifetime time.Duration) (cookie models.Cookie, err error) {
	cookie = s.createCookie(userId, lifetime*time.Second)

	err = s.authorizationStorage.InsertCookie(cookie)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}
	return
}


func (s *service) GetUserIdByCookie(cookieValue string) (userID int, err error) {
	cookie, err := s.authorizationStorage.SelectCookieByValue(cookieValue)
	if err != nil {
		if err == sql.ErrNoRows {
			err = s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorNotValidCookie, err)
			return
		}
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	if cookie.Expiration.Before(time.Now()) {
		err = s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorNotValidCookie, err)
		return
	}

	return cookie.UserID, err
}

func (s *service) createCookie(userID int, expiration time.Duration) (cookie models.Cookie) {
	cookie.UserID = userID
	cookie.Expiration = time.Now().Add(expiration)
	cookie.Value = common.GetRandomString(int64(config.Config.SessionCookieLength),
		common.Capitals, common.Lowers, common.Numbers, common.LightSymbols)
	return
}

func (s *service) checkLogoPass(userID int, userPassword string) (err error) {
	basePass, err := s.authorizationStorage.SelectPasswordByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorUserNotFound, ErrorUserNotFound)
		}
		return s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	hashPassword := s.hashPassword([]byte(userPassword), basePass.Salt)
	if !bytes.Equal(hashPassword, basePass.HashPassword) {
		return s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorIncorrectPassword, ErrorIncorrectPassword)
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

func (s *service) createUser(user models.Registration) (userID int, err error) {
	userID, err = s.accountClient.CreateUser(models2.User{
		Email:     user.Login,
		Name:      user.Name,
		Surname:   user.Surname,
		AvatarURL: user.Avatar,
	})
	return
}

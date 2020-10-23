package authorizationStorage

import (
	"database/sql"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
)

type Storage interface {
	InsertPassword(pass models.Password) (err error)
	UpdatePassword(pass models.Password) (err error)
	SelectPasswordByUserID(userID int) (pass models.Password, err error)

	InsertCookie(cookie models.Cookie) (err error)
	SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertCookie(cookie models.Cookie) (err error) {
	panic("implement me")
}

func (s *storage) SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error) {
	panic("implement me")
}

func (s *storage) InsertPassword(pass models.Password) (err error) {
	panic("implement me")
}

func (s *storage) UpdatePassword(pass models.Password) (err error) {
	panic("implement me")
}

func (s *storage) SelectPasswordByUserID(userID int) (pass models.Password, err error) {
	panic("implement me")
}

package authorizationStorage

import (
	"database/sql"
	"fmt"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
)

type Storage interface {
	InsertPassword(pass models.Password) (err error)
	UpdatePassword(pass models.Password) (err error)
	SelectPasswordByUserID(userID int) (pass models.Password, err error)

	InsertCookie(cookie models.Cookie) (err error)
	SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error)
}

const (
	sessionsTable = "sessions"
	passwordsTable = "passwords"
)

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertCookie(cookie models.Cookie) (err error) {
	const tmplQuery = `
	INSERT INTO %s 
	(user_id, cookie, expiration)
	VALUES ($1, $2, $3)`

	query := fmt.Sprintf(tmplQuery, sessionsTable)
	res, err := s.db.Exec(query, cookie.UserID, cookie.Value, cookie.Expiration)
	if err != nil {
		return
	}
	c, err := res.RowsAffected()
	if err != nil {
		return
	}
	if c != 1 {
		err = fmt.Errorf("rows affected: wrong result")
	}
	return
}

func (s *storage) SelectCookieByValue(cookieValue string) (cookie models.Cookie, err error) {
	const tmplQuery = `SELECT user_id, expiration FROM %s WHERE cookie=$1`

	query := fmt.Sprintf(tmplQuery, sessionsTable)
	cookie.Value = cookieValue
	err = s.db.QueryRow(query, cookieValue).Scan(&cookie.UserID, &cookie.Expiration)

	return
}

func (s *storage) InsertPassword(pass models.Password) (err error) {
	const tmplQuery = `INSERT INTO %s as p (user_id, password_hash, salt, update_at) VALUES($1, $2, $3, $4) ON CONFLICT(user_id) DO UPDATE SET password_hash=$2, salt=$3, update_at=$4 WHERE p.user_id=$1`
	const tmplQueryFirst = `INSERT INTO %s (password_hash, salt, update_at) VALUES($1, $2, $3)`
	if pass.UserID == 0 {
		query := fmt.Sprintf(tmplQueryFirst, passwordsTable)
		_, err = s.db.Exec(query, pass.HashPassword, pass.Salt, pass.UpdateAt)
		return
	}
	query := fmt.Sprintf(tmplQuery, passwordsTable)
	_, err = s.db.Exec(query, pass.UserID, pass.HashPassword, pass.Salt, pass.UpdateAt)
	return
}

func (s *storage) UpdatePassword(pass models.Password) (err error) {
	const tmplQuery = `UPDATE %s SET password_hash=$2, salt=$3, update_at=$4 WHERE user_id = $1`

	query := fmt.Sprintf(tmplQuery, passwordsTable)
	res, err := s.db.Exec(query, pass.UserID, pass.HashPassword, pass.Salt, pass.UpdateAt)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		err = fmt.Errorf("bad rows count affected")
	}
	return
}

func (s *storage) SelectPasswordByUserID(userID int) (pass models.Password, err error) {
	const tmplQuery = `SELECT password_hash, salt, update_at FROM %s WHERE user_id=$1`

	query := fmt.Sprintf(tmplQuery, passwordsTable)

	err = s.db.QueryRow(query, userID).Scan(&pass.HashPassword, &pass.Salt, &pass.UpdateAt)

	return
}

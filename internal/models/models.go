package models

import "time"

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Uid 	 int 	`json:'uid'`
}

type Registration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Cookie struct {
	UserID     int
	Value      string
	Expiration time.Time
}

type Password struct {
	UserID   int
	HashPassword []byte
	Salt     []byte
	UpdateAt time.Time
}

type CheckAuthRequest struct {
	SessionToken string	`json:"cookie"`
}

type CheckAuthResponse struct {
	Uid int	`json:"uid"`
	//Email string	`json:"email"`
}

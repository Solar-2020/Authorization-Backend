package models

import "time"

type Authorization struct {
	Login    string `json:"login" validation:"required"`
	Password string `json:"password" validation:"required"`
}

type AuthorizationResponse struct {
	Login  string `json:"login"`
	Status string `json:"status"`
	Uid    int    `json:'uid'`
}

type Registration struct {
	Login    string `json:"login" validation:"required,email"`
	Password string `json:"password,omitempty" validation:"required"`
	Avatar   string `json:"avatar"`
	Name     string `json:"name" validation:"required"`
	Surname  string `json:"surname"`
}

type RegistrationResponse struct {
	Registration
	Uid int `json:'uid'`
}

type Cookie struct {
	UserID     int
	Value      string
	Expiration time.Time
}

type Password struct {
	UserID       int
	HashPassword []byte
	Salt         []byte
	UpdateAt     time.Time
}

type CheckAuthRequest struct {
	SessionToken string `json:"cookie" validation"required"`
}

type CheckAuthResponse struct {
	Uid int `json:"uid"`
	//Email string	`json:"email"`
}

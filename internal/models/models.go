package models

import "time"

type Authorization struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
	HashPassword string
	Salt     string
	UpdateAt time.Time
}

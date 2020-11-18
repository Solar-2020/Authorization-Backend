package client

import "github.com/pkg/errors"

var (
	ErrorUnknownStatusCode = "Unknown status code %v"
	ErrorInvalidSecretKey  = errors.New("Invalid secret key")
)

type GetUserIdByCookieResponse struct {
	UserID int `json:"uid"`
}

type httpError struct {
	Error string `json:"error"`
}

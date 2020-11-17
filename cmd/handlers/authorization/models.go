package authorizationHandler

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type authorizationService interface {
	Authorization(request models.Authorization) (cookie models.Cookie, err error)
	Registration(request models.Registration) (cookie models.Cookie, err error)
	Yandex(userToken string) (cookie models.Cookie, err error)
	GetUserIdByCookie(cookieValue string) (userID int, err error)
}

type authorizationTransport interface {
	AuthorizationDecode(ctx *fasthttp.RequestCtx) (request models.Authorization, err error)
	AuthorizationEncode(ctx *fasthttp.RequestCtx, resp models.AuthorizationResponse, cookie models.Cookie) (err error)

	RegistrationDecode(ctx *fasthttp.RequestCtx) (request models.Registration, err error)
	RegistrationEncode(ctx *fasthttp.RequestCtx, resp models.RegistrationResponse, cookie models.Cookie) (err error)

	GetUserIdByCookieDecode(ctx *fasthttp.RequestCtx) (request models.CheckAuthRequest, err error)
	GetUserIdByCookieEncode(ctx *fasthttp.RequestCtx, response models.CheckAuthResponse) (err error)

	YandexDecode(ctx *fasthttp.RequestCtx) (userToken string, err error)
	YandexEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error)
	NewError(httpCode int, responseError error, fullError error) (err error)
}

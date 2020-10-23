package authorizationHandler

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type authorizationService interface {
	Authorization(request models.Authorization) (cookie models.Cookie, err error)
	Registration(request models.Registration) (cookie models.Cookie, err error)
	GetUserIdByCookie(cookieValue string) (userID int, err error)
}

type authorizationTransport interface {
	AuthorizationDecode(ctx *fasthttp.RequestCtx) (request models.Authorization, err error)
	AuthorizationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error)

	RegistrationDecode(ctx *fasthttp.RequestCtx) (request models.Registration, err error)
	RegistrationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error)

	GetUserIdByCookieDecode(ctx *fasthttp.RequestCtx) (cookieValue string, err error)
	GetUserIdByCookieEncode(ctx *fasthttp.RequestCtx, userID int) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error) (err error)
	ServeFatalError(ctx *fasthttp.RequestCtx)
}

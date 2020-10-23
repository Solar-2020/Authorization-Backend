package authorization

import (
	"encoding/json"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Transport interface {
	AuthorizationDecode(ctx *fasthttp.RequestCtx) (request models.Authorization, err error)
	AuthorizationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error)

	RegistrationDecode(ctx *fasthttp.RequestCtx) (request models.Registration, err error)
	RegistrationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error)

	GetUserIdByCookieDecode(ctx *fasthttp.RequestCtx) (cookieValue string, err error)
	GetUserIdByCookieEncode(ctx *fasthttp.RequestCtx, userID int) (err error)
}

type transport struct {
}

func NewTransport() Transport {
	return &transport{}
}

func (t transport) AuthorizationDecode(ctx *fasthttp.RequestCtx) (request models.Authorization, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)

	return
}

func (t transport) AuthorizationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	panic("implement me")
	return
}

func (t transport) RegistrationDecode(ctx *fasthttp.RequestCtx) (request models.Registration, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)

	return
}

func (t transport) RegistrationEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	panic("implement me")
	return
}

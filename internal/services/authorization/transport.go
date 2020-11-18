package authorization

import (
	"encoding/json"
	"github.com/Solar-2020/Authorization-Backend/cmd/config"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/go-playground/validator"
	"github.com/valyala/fasthttp"
)

type transport struct {
	validator *validator.Validate
}

func NewTransport() *transport {
	return &transport{
		validator: validator.New(),
	}
}

func (t transport) AuthorizationDecode(ctx *fasthttp.RequestCtx) (request models.Authorization, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) AuthorizationEncode(ctx *fasthttp.RequestCtx, resp models.AuthorizationResponse, cookie models.Cookie) (err error) {
	body, err := json.Marshal(resp)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	t.setCookie(ctx, cookie)
	return
}

func (t transport) RegistrationDecode(ctx *fasthttp.RequestCtx) (request models.Registration, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) RegistrationEncode(ctx *fasthttp.RequestCtx, resp models.RegistrationResponse, cookie models.Cookie) (err error) {
	resp.Password = ""
	body, err := json.Marshal(resp)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	t.setCookie(ctx, cookie)
	return
}

func (t transport) GetUserIdByCookieDecode(ctx *fasthttp.RequestCtx) (request models.CheckAuthRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		return
	}
	err = t.validator.Struct(request)
	return
}

func (t transport) GetUserIdByCookieDecodeV2(ctx *fasthttp.RequestCtx) (request models.CheckAuthRequest, err error) {
	request.SessionToken = string(ctx.QueryArgs().Peek("session_cookie"))
	return
}

func (t transport) GetUserIdByCookieEncode(ctx *fasthttp.RequestCtx, response models.CheckAuthResponse) (err error) {
	body, err := json.Marshal(response)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) YandexDecode(ctx *fasthttp.RequestCtx) (userToken string, err error) {
	userToken = ctx.UserValue("userToken").(string)
	return
}

func (t transport) YandexEncode(ctx *fasthttp.RequestCtx, cookie models.Cookie) (err error) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	t.setCookie(ctx, cookie)
	return
}

func (t *transport) setCookie(ctx *fasthttp.RequestCtx, src models.Cookie) {
	cookie := fasthttp.Cookie{}
	cookie.SetKey(config.Config.SessionCookieName)
	cookie.SetDomain(string(ctx.Request.Host()))
	cookie.SetPath("/")
	cookie.SetValue(src.Value)
	cookie.SetExpire(src.Expiration)
	ctx.Response.Header.SetCookie(&cookie)
}

package authorizationHandler

import (
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Authorization(ctx *fasthttp.RequestCtx)
	Registration(ctx *fasthttp.RequestCtx)
	Yandex(ctx *fasthttp.RequestCtx)
	GetUserIdByCookie(ctx *fasthttp.RequestCtx)
	GetUserIdByCookieV2(ctx *fasthttp.RequestCtx)
}

type handler struct {
	authorizationService   authorizationService
	authorizationTransport authorizationTransport
	errorWorker            errorWorker
}

func NewHandler(authorizationService authorizationService, authorizationTransport authorizationTransport, errorWorker errorWorker) Handler {
	return &handler{
		authorizationService:   authorizationService,
		authorizationTransport: authorizationTransport,
		errorWorker:            errorWorker,
	}
}

func (h *handler) Authorization(ctx *fasthttp.RequestCtx) {
	auth, err := h.authorizationTransport.AuthorizationDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	cookie, err := h.authorizationService.Authorization(auth)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	resp := models.AuthorizationResponse{
		Login:  auth.Login,
		Status: "OK",
		Uid:    cookie.UserID,
	}

	err = h.authorizationTransport.AuthorizationEncode(ctx, resp, cookie)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Registration(ctx *fasthttp.RequestCtx) {
	auth, err := h.authorizationTransport.RegistrationDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	cookie, err := h.authorizationService.Registration(auth)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	resp := models.RegistrationResponse{
		Registration: auth,
		Uid:          cookie.UserID,
	}
	err = h.authorizationTransport.RegistrationEncode(ctx, resp, cookie)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Yandex(ctx *fasthttp.RequestCtx) {
	auth, err := h.authorizationTransport.YandexDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	cookie, err := h.authorizationService.Yandex(auth)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.authorizationTransport.YandexEncode(ctx, cookie)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) GetUserIdByCookie(ctx *fasthttp.RequestCtx) {
	req, err := h.authorizationTransport.GetUserIdByCookieDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	userID, err := h.authorizationService.GetUserIdByCookie(req.SessionToken)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.authorizationTransport.GetUserIdByCookieEncode(ctx, models.CheckAuthResponse{Uid: userID})
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) GetUserIdByCookieV2(ctx *fasthttp.RequestCtx) {
	req, err := h.authorizationTransport.GetUserIdByCookieDecodeV2(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	userID, err := h.authorizationService.GetUserIdByCookie(req.SessionToken)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.authorizationTransport.GetUserIdByCookieEncode(ctx, models.CheckAuthResponse{Uid: userID})
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

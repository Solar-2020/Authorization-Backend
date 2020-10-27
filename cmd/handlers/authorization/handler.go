package authorizationHandler

import (
	"fmt"
	"github.com/Solar-2020/Authorization-Backend/internal/models"
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Authorization(ctx *fasthttp.RequestCtx)
	Registration(ctx *fasthttp.RequestCtx)
	GetUserIdByCookie(ctx *fasthttp.RequestCtx)
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
	fmt.Println("New incoming request: POST /authorization/authorization")
	auth, err := h.authorizationTransport.AuthorizationDecode(ctx)
	if err != nil {
		fmt.Println("Create: cannot decode request")
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	cookie, err := h.authorizationService.Authorization(auth)
	if err != nil {
		fmt.Println("Create: bad usecase: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.authorizationTransport.AuthorizationEncode(ctx, cookie)
	if err != nil {
		fmt.Println("Create: cannot encode response: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) Registration(ctx *fasthttp.RequestCtx) {
	fmt.Println("New incoming request: POST /authorization/authorization")
	auth, err := h.authorizationTransport.RegistrationDecode(ctx)
	if err != nil {
		fmt.Println("Create: cannot decode request")
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	cookie, err := h.authorizationService.Registration(auth)
	if err != nil {
		fmt.Println("Create: bad usecase: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}

	err = h.authorizationTransport.RegistrationEncode(ctx, cookie)
	if err != nil {
		fmt.Println("Create: cannot encode response: ", err)
		err = h.errorWorker.ServeJSONError(ctx, err)
		if err != nil {
			h.errorWorker.ServeFatalError(ctx)
		}
		return
	}
}

func (h *handler) GetUserIdByCookie(ctx *fasthttp.RequestCtx) {
	req, err := h.authorizationTransport.GetUserIdByCookieDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	userID, err := h.authorizationService.GetUserIdByCookie(req.SessionToken)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.authorizationTransport.GetUserIdByCookieEncode(ctx, models.CheckAuthResponse{Uid: userID})
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) handleError(err error, ctx *fasthttp.RequestCtx) {
	err = h.errorWorker.ServeJSONError(ctx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx)
	}
	return
}
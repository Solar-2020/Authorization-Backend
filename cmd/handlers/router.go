package handlers

import (
	"fmt"
	authorizationHandler "github.com/Solar-2020/Authorization-Backend/cmd/handlers/authorization"
	"github.com/Solar-2020/Authorization-Backend/internal/errorWorker"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

func NewFastHttpRouter(authorization authorizationHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	//router.Handle("GET", "/health", check)

	router.PanicHandler = panicHandler

	router.Handle("POST", "/authorization/authorization", middleware.CORS(authorization.Authorization))
	router.Handle("GET", "/authorization/registration", middleware.CORS(authorization.Registration))

	router.Handle("GET", "/authorization/user-id", middleware.CORS(authorization.GetUserIdByCookie))



	return router
}

func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	fmt.Printf("Request falied with panic: %s, error: %v\nTrace:\n", string(ctx.Request.RequestURI()), err)
	fmt.Println(string(debug.Stack()))
	errorWorker.NewErrorWorker().ServeFatalError(ctx)
}

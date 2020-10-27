package handlers

import (
	authorizationHandler "github.com/Solar-2020/Authorization-Backend/cmd/handlers/authorization"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(authorization authorizationHandler.Handler, middleware httputils.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.PanicHandler = httputils.PanicHandler
	middlewareChain := httputils.NewLogCorsChain(middleware)

	router.Handle("POST", "/auth/login", middlewareChain(authorization.Authorization))
	router.Handle("PUT", "/auth/signup", middlewareChain(authorization.Registration))

	router.Handle("POST", "/auth/cookie", middlewareChain(authorization.GetUserIdByCookie))

	return router
}
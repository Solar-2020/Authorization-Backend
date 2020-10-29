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

	router.Handle("POST", "/api/auth/login", middlewareChain(authorization.Authorization))
	router.Handle("PUT", "/api/auth/signup", middlewareChain(authorization.Registration))

	router.Handle("POST", "/api/auth/cookie", middlewareChain(authorization.GetUserIdByCookie))

	return router
}
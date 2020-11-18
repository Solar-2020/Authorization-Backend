package main

import (
	"database/sql"
	account "github.com/Solar-2020/Account-Backend/pkg/client"
	"github.com/Solar-2020/Authorization-Backend/cmd/config"
	"github.com/Solar-2020/Authorization-Backend/cmd/handlers"
	authorizationHandler "github.com/Solar-2020/Authorization-Backend/cmd/handlers/authorization"

	"github.com/Solar-2020/Authorization-Backend/internal/services/authorization"
	"github.com/Solar-2020/Authorization-Backend/internal/storages/authorizationStorage"
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	err := envconfig.Process("", &config.Config)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	authorizationDB, err := sql.Open("postgres", config.Config.AuthorizationDataBaseConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	authorizationDB.SetMaxIdleConns(5)
	authorizationDB.SetMaxOpenConns(10)

	errorWorker := errorWorker.NewErrorWorker()

	authorizationStorage := authorizationStorage.NewStorage(authorizationDB)
	accountClient := account.NewClient(config.Config.AccountServiceHost, config.Config.ServerSecret)
	authorizationService := authorization.NewService(authorizationStorage, accountClient, errorWorker)
	authorizationTransport := authorization.NewTransport()

	authorizationHandler := authorizationHandler.NewHandler(authorizationService, authorizationTransport, errorWorker)

	middlewares := httputils.NewMiddleware()

	server := fasthttp.Server{
		Handler: handlers.NewFastHttpRouter(authorizationHandler, middlewares).Handler,
	}

	go func() {
		log.Info().Str("msg", "start server").Str("port", config.Config.Port).Send()
		if err := server.ListenAndServe(":" + config.Config.Port); err != nil {
			log.Error().Str("msg", "server run failure").Err(err).Send()
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	defer func(sig os.Signal) {

		log.Info().Str("msg", "received signal, exiting").Str("signal", sig.String()).Send()

		if err := server.Shutdown(); err != nil {
			log.Error().Str("msg", "server shutdown failure").Err(err).Send()
		}

		//dbConnection.Shutdown()
		log.Info().Str("msg", "goodbye").Send()
	}(<-c)
}

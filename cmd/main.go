package main

import (
	"database/sql"
	"github.com/Solar-2020/Authorization-Backend/cmd/handlers"
	authorizationHandler "github.com/Solar-2020/Authorization-Backend/cmd/handlers/authorization"
	"github.com/Solar-2020/Authorization-Backend/internal/errorWorker"
	"github.com/Solar-2020/Authorization-Backend/internal/services/authorization"
	"github.com/Solar-2020/Authorization-Backend/internal/storages/authorizationStorage"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	Port                          string `envconfig:"PORT" default:"8099"`
	AuthorizationDataBaseConnectionString string `envconfig:"AUTHORIZATION_DB_CONNECTION_STRING" default:"-"`
}

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	authorizationDB, err := sql.Open("postgres", cfg.AuthorizationDataBaseConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	authorizationDB.SetMaxIdleConns(5)
	authorizationDB.SetMaxOpenConns(10)

	errorWorker := errorWorker.NewErrorWorker()

	authorizationStorage := authorizationStorage.NewStorage(authorizationDB)
	authorizationService := authorization.NewService(authorizationStorage)
	authorizationTransport := authorization.NewTransport()

	authorizationHandler := authorizationHandler.NewHandler(authorizationService, authorizationTransport, errorWorker)

	middlewares := handlers.NewMiddleware()

	server := fasthttp.Server{
		Handler: handlers.NewFastHttpRouter(authorizationHandler, middlewares).Handler,
	}

	go func() {
		log.Info().Str("msg", "start server").Str("port", cfg.Port).Send()
		if err := server.ListenAndServe(":" + cfg.Port); err != nil {
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

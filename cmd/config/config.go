package config

import "github.com/Solar-2020/GoUtils/common"

var (
	Config config
)

type config struct {
	common.SharedConfig
	AuthorizationDataBaseConnectionString 	string `envconfig:"AUTHORIZATION_DB_CONNECTION_STRING" default:"-"`
	DefaultCookieLifetime					int64	`envconfig:"DEFAULT_COOKIE_LIFETIME" default:"100000"`
	SessionCookieLength						int 	`envconfig:"SESSING_COOKIE_LENGTH" default:"10"`
	SessionCookieName						string	`envconfig:"SESSING_COOKIE_NAME" default:"SessionToken"`
}

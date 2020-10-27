package config

var (
	Config config
)

type config struct {
	Port                          			string `envconfig:"PORT" default:"8099"`
	AuthorizationDataBaseConnectionString 	string `envconfig:"AUTHORIZATION_DB_CONNECTION_STRING" default:"-"`
	DefaultCookieLifetime					int64	`envconfig:"DEFAULT_COOKIE_LIFETIME" default:"100000"`
	SessionCookieLength						int 	`envconfig:"SESSING_COOKIE_LENGTH" default:"10"`
}

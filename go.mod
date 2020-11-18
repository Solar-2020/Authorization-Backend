module github.com/Solar-2020/Authorization-Backend

go 1.14

require (
	github.com/Solar-2020/Account-Backend v1.0.5
	github.com/Solar-2020/GoUtils v1.0.3
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.8.0
	github.com/pkg/errors v0.8.1
	github.com/rs/zerolog v1.20.0
	github.com/valyala/fasthttp v1.16.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
)

// replace github.com/Solar-2020/GoUtils => ../GoUtils

// replace github.com/Solar-2020/Account-Backend => ../Account-Backend

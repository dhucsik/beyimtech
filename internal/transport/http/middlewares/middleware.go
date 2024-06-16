package middlewares

import "github.com/valyala/fasthttp"

type Middleware interface {
	Handler(fasthttp.RequestHeader) fasthttp.RequestHandler
}

package middleware

import (
	"net/http"

	"github.com/handson-go/chap6/complex_server/config"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicMiddleware(mux, c), c)
}

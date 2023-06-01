package handlers

import (
	"net/http"

	"github.com/handson-go/chap6/complex_server/config"
)

func Register(mux *http.ServeMux, conf config.AppConfig) {
	mux.Handle("/healthz", &app{conf: conf, handler: healthzHandler})
	mux.Handle("/api", &app{conf: conf, handler: apiHandler})
	mux.Handle("/panic", &app{conf: conf, handler: panicHandler})
}

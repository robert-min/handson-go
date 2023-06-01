package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/handson-go/chap6/complex_server/config"
	"github.com/handson-go/chap6/complex_server/handlers"
	"github.com/handson-go/chap6/complex_server/middleware"
)

func setupServer(mux *http.ServeMux, w io.Writer) http.Handler {
	conf := config.InitConfig(w)

	handlers.Register(mux, conf)
	return middleware.RegisterMiddleware(mux, conf)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	wrappedMux := setupServer(mux, os.Stdout)

	log.Fatal(http.ListenAndServe(listenAddr, wrappedMux))
}

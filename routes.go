package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/wschat/handlers"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndPoint))

	fileSrv := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileSrv))
	return mux
}

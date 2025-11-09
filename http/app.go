package main

import (
	"fmt"
	"net/http"
	"playgoround/http/internal/config"
	"playgoround/http/internal/handler"
	"strconv"
)

type App struct {
	config   *config.Config
	handlers []handler.HttpHandler
}

func New(cfg *config.Config) *App {
	handlers := make([]handler.HttpHandler, 0)

	monitorHandler := handler.NewMonitorHandler(
		handler.WithHealthCheckResponse("foo"),
	)
	handlers = append(handlers, monitorHandler)

	return &App{
		config:   cfg,
		handlers: handlers,
	}
}

func (a *App) Start() {
	mux := http.NewServeMux()
	for _, httpHandler := range a.handlers {
		httpHandler.Register(mux)
	}

	err := http.ListenAndServe(":"+strconv.Itoa(a.config.HTTP.Port), mux)
	if err != nil {
		fmt.Println("failed to start server", err)
		panic(err)
	}
}

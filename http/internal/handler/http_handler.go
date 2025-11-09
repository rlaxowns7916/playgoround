package handler

import "net/http"

type HttpHandler interface {
	Register(mux *http.ServeMux)
}

type Option interface {
	apply(h HttpHandler)
}

type optionFunc func(h HttpHandler)

func (o optionFunc) apply(h HttpHandler) {
	o(h)
}

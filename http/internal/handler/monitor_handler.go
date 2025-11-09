package handler

import "net/http"

type MonitorHandler struct {
	healthCheckResponse string
}

func NewMonitorHandler(options ...Option) *MonitorHandler {
	handler := &MonitorHandler{
		healthCheckResponse: "ok",
	}
	for _, opt := range options {
		opt.apply(handler)
	}

	return handler
}

func WithHealthCheckResponse(r string) Option {
	return optionFunc(func(h HttpHandler) {
		if mh, ok := h.(*MonitorHandler); ok {
			mh.healthCheckResponse = r
		}
	})
}

func (u *MonitorHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", u.health)
}

func (u *MonitorHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(u.healthCheckResponse))
}

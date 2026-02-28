package server

import (
	"net/http"

	"github.com/AndreiBerezin/pixoo64/pkg/log"
)

type MicHandler interface {
	SetOnAir(on bool)
}

type Server struct {
	address string
	handler MicHandler
}

func New(addr string, handler MicHandler) *Server {
	if addr == "" {
		addr = ":8080"
	}
	return &Server{address: addr, handler: handler}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /mic/on", func(w http.ResponseWriter, r *http.Request) {
		s.handler.SetOnAir(true)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /mic/off", func(w http.ResponseWriter, r *http.Request) {
		s.handler.SetOnAir(false)
		w.WriteHeader(http.StatusOK)
	})

	log.Info("http server listening on " + s.address)
	if err := http.ListenAndServe(s.address, mux); err != nil {
		log.Fatal("http server error: " + err.Error())
	}
}

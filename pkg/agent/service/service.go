package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	router *mux.Router
}

func NewService() *Service {
	s := &Service{
		router: mux.NewRouter(),
	}

	// This API will be filled out more later
	// Just leaving a stub here for now
	s.router.HandleFunc("/", func(http.ResponseWriter, *http.Request) {}).Methods("POST")

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

package service

import (
	"net/http"

	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/gorilla/mux"
)

type Service struct {
	variables variables.Interface
	confDir   string
	router    *mux.Router
}

func NewService(variables variables.Interface, confDir string) *Service {
	s := &Service{
		variables: variables,
		confDir:   confDir,
		router:    mux.NewRouter(),
	}

	s.router.HandleFunc("/ssh", s.ssh).Methods("POST")
	s.router.HandleFunc("/execute", s.execute).Methods("POST")

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

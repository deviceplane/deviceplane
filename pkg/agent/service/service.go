package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"sync"

	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/gliderlabs/ssh"
	"github.com/gorilla/mux"
	gossh "golang.org/x/crypto/ssh"

	"net/http/pprof"
)

type Service struct {
	variables variables.Interface
	confDir   string
	router    *mux.Router

	signer     ssh.Signer
	signerLock sync.Mutex
}

func NewService(
	variables variables.Interface, confDir string,
) *Service {
	s := &Service{
		variables: variables,
		confDir:   confDir,
		router:    mux.NewRouter(),
	}
	//go s.getSigner()

	s.router.HandleFunc("/ssh", s.ssh).Methods("POST")
	s.router.HandleFunc("/reboot", s.reboot).Methods("POST")

	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.router.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)

	return s
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) getSigner() (ssh.Signer, error) {
	s.signerLock.Lock()
	defer s.signerLock.Unlock()

	if s.signer != nil {
		return s.signer, nil
	}

	hostSignerKey := s.variables.GetHostSignerKey()

	var key *rsa.PrivateKey
	var err error
	if hostSignerKey == "" {
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode([]byte(hostSignerKey))
		if block == nil {
			return nil, err
		}
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	signer, err := gossh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}

	s.signer = signer

	return s.signer, nil
}

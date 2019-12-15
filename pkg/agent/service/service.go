package service

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"sync"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/metrics"
	"github.com/deviceplane/deviceplane/pkg/agent/netns"
	"github.com/deviceplane/deviceplane/pkg/agent/supervisor"
	"github.com/deviceplane/deviceplane/pkg/agent/variables"
	"github.com/deviceplane/deviceplane/pkg/engine"
	"github.com/gliderlabs/ssh"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	gossh "golang.org/x/crypto/ssh"
)

const (
	HostMetricsBrokenMsg = "host metrics are not working, check agent logs for details"
)

var (
	hostMetricsFallbackHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, HostMetricsBrokenMsg, http.StatusInternalServerError)
	}))
)

type Service struct {
	variables        variables.Interface
	supervisorLookup supervisor.Lookup
	confDir          string
	netnsManager     *netns.Manager
	router           *mux.Router

	signer     ssh.Signer
	signerLock sync.Mutex
}

func NewService(
	variables variables.Interface, supervisorLookup supervisor.Lookup,
	engine engine.Engine, confDir string,
) *Service {
	s := &Service{
		variables:        variables,
		supervisorLookup: supervisorLookup,
		confDir:          confDir,
		netnsManager:     netns.NewManager(engine),
		router:           mux.NewRouter(),
	}
	go s.getSigner()

	hostMetricsHandler, err := metrics.HostMetricsHandler(nil)
	if err != nil {
		log.WithError(err).Error("create host metrics handler")
		hostMetricsHandler = &hostMetricsFallbackHandler
	}

	s.router.HandleFunc("/ssh", s.ssh).Methods("POST")
	s.router.HandleFunc("/reboot", s.reboot).Methods("POST")
	s.router.HandleFunc("/applications/{application}/services/{service}/imagepullprogress", s.imagePullProgress).Methods("GET")
	s.router.HandleFunc("/applications/{application}/services/{service}/metrics", s.metrics).Methods("GET")
	s.router.Handle("/metrics/host", *hostMetricsHandler)
	s.router.Handle("/metrics/agent", promhttp.Handler())

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

	// Generate
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	signer, err := gossh.NewSignerFromKey(key)
	if err != nil {
		return nil, err
	}

	s.signer = signer

	return s.signer, nil
}

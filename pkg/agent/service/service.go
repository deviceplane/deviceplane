package service

import (
	"bytes"
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
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/filtering"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/gliderlabs/ssh"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	gossh "golang.org/x/crypto/ssh"
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

	s.router.HandleFunc("/ssh", s.ssh).Methods("POST")
	s.router.HandleFunc("/reboot", s.reboot).Methods("POST")
	s.router.HandleFunc("/applications/{application}/services/{service}/imagepullprogress", s.imagePullProgress).Methods("GET")
	s.router.HandleFunc("/applications/{application}/services/{service}/metrics", s.metrics).Methods("GET")
	s.router.Handle("/metrics/host", newHostMetricsHandler())
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

func newHostMetricsHandler() http.Handler {
	var hostMetricsHandler http.Handler
	unfilteredHostMetricsHandler, err := metrics.HostMetricsHandler(nil)
	if err == nil { // Proxy handler response and filter node prefix
		hostMetricsHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			rwHttp := utils.ResponseWriter{
				Headers: w.Header(),
				Writer:  &buf,
			}

			(*unfilteredHostMetricsHandler).ServeHTTP(&rwHttp, r)

			w.WriteHeader(rwHttp.Status)
			if rwHttp.Status == http.StatusOK {
				rawHostMetricsString := buf.String()
				filteredHostMetrics := filtering.FilterNodePrefix(rawHostMetricsString)
				w.Write([]byte(filteredHostMetrics))
			} else {
				w.Write(buf.Bytes())
			}
		}))
	} else {
		log.WithError(err).Error("create host metrics handler")
		hostMetricsHandler = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "host metrics are not working, check agent logs for details", http.StatusInternalServerError)
		}))
	}
	return hostMetricsHandler
}

package metrics

import (
	"bytes"
	"context"
	"net/http"
	"sync"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/agent/netns"
	"github.com/deviceplane/deviceplane/pkg/agent/supervisor"
	"github.com/deviceplane/deviceplane/pkg/metrics/datadog/processing"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/pkg/errors"
)

var once sync.Once
var hostMetricsHandler http.Handler

func GetFilteredHostMetrics(ctx context.Context) (*bytes.Buffer, error) {
	h := FilteredHostMetricsHandler()

	var buf bytes.Buffer
	rwHttp := utils.ResponseWriter{
		Writer: &buf,
	}

	r, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		return nil, err
	}
	(h).ServeHTTP(&rwHttp, r)

	return &buf, nil
}

func FilteredHostMetricsHandler() http.Handler {
	once.Do(func() {
		unfilteredHostMetricsHandler, err := HostMetricsHandler(nil)
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
					filteredHostMetrics := processing.PrefilterNodePrefix(rawHostMetricsString)
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
	})

	return hostMetricsHandler
}

type ServiceMetricsFetcher struct {
	supervisorLookup supervisor.Lookup
	netnsManager     *netns.Manager
}

func NewServiceMetricsFetcher(
	supervisorLookup supervisor.Lookup,
	netnsManager *netns.Manager,
) *ServiceMetricsFetcher {
	return &ServiceMetricsFetcher{
		supervisorLookup: supervisorLookup,
		netnsManager:     netnsManager,
	}
}

func (s *ServiceMetricsFetcher) ContainerServiceMetrics(ctx context.Context, applicationID, service string, port int, path string) (*http.Response, error) {
	containerID, ok := s.supervisorLookup.GetContainerID(applicationID, service)
	if !ok {
		return nil, errors.New("could not get container ID")
	}

	resp, err := s.netnsManager.ProcessRequest(
		ctx, containerID, port, string(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "container could not process request")
	}

	return resp, nil
}

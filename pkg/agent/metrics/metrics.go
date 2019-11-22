package metrics

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/node_exporter/collector"
)

// This should be more than enough concurrent requests. We're only polling this
// once per device (non-concurrently), so other requests come from customer API
// calls. If we exceed this count, we're probably doing something wrong.
const MaxRequestsInFlight = 3

func HostMetricsHandler() (*http.Handler, error) {
	nc, err := newNodeCollector(
		"cpu",
		"diskstats",
		"filesystem",
		"loadavg",
		"meminfo",
		"textfile",
		"time",
		"netdev",
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	if err := r.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
	}

	metricsRegistry := prometheus.NewRegistry()
	handler := promhttp.InstrumentMetricHandler(
		metricsRegistry,
		promhttp.HandlerFor(
			prometheus.Gatherers{metricsRegistry, r},
			promhttp.HandlerOpts{
				ErrorLog:            log.NewErrorLogger(),
				ErrorHandling:       promhttp.ContinueOnError,
				MaxRequestsInFlight: MaxRequestsInFlight,
				Registry:            metricsRegistry,
			},
		),
	)
	return &handler, nil
}

var collectorCreators = map[string]func() (collector.Collector, error){
	"cpu":         collector.NewCPUCollector,
	"diskstats":   collector.NewDiskstatsCollector,
	"filesystem":  collector.NewFilesystemCollector,
	"loadavg":     collector.NewLoadavgCollector,
	"meminfo":     collector.NewMeminfoCollector,
	"textfile":    collector.NewTextFileCollector,
	"time":        collector.NewTimeCollector,
	"runit":       collector.NewRunitCollector,
	"supervisord": collector.NewSupervisordCollector,
	"netdev":      collector.NewNetDevCollector,
	"ntp":         collector.NewNtpCollector,
}

func newNodeCollector(collectorNames ...string) (*collector.NodeCollector, error) {
	collectors := make(map[string]collector.Collector)

	for _, name := range collectorNames {
		cCreator, exists := collectorCreators[name]
		if !exists {
			return nil, errors.New("invalid collector")
		}

		c, err := cCreator()
		if err != nil {
			return nil, err
		}

		collectors[name] = c
	}
	return &collector.NodeCollector{Collectors: collectors}, nil
}
